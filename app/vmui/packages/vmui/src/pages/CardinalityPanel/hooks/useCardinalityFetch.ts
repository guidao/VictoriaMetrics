import { ErrorTypes } from "../../../types";
import { useAppState } from "../../../state/common/StateContext";
import { useEffect, useState } from "preact/compat";
import { CardinalityRequestsParams, getCardinalityInfo } from "../../../api/tsdb";
import { TSDBStatus } from "../types";
import AppConfigurator from "../appConfigurator";
import { useSearchParams } from "react-router-dom";
import dayjs from "dayjs";
import { DATE_FORMAT } from "../../../constants/date";

export const useFetchQuery = (): {
  fetchUrl?: string[],
  isLoading: boolean,
  error?: ErrorTypes | string
  appConfigurator: AppConfigurator,
} => {
  const appConfigurator = new AppConfigurator();

  const [searchParams] = useSearchParams();
  const match = searchParams.get("match");
  const focusLabel = searchParams.get("focusLabel");
  const topN = +(searchParams.get("topN") || 10);
  const date = searchParams.get("date") || dayjs().tz().format(DATE_FORMAT);

  const { serverUrl } = useAppState();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<ErrorTypes | string>();
  const [tsdbStatus, setTSDBStatus] = useState<TSDBStatus>(appConfigurator.defaultTSDBStatus);

  const fetchCardinalityInfo = async (requestParams: CardinalityRequestsParams) => {
    if (!serverUrl) return;
    setError("");
    setIsLoading(true);
    setTSDBStatus(appConfigurator.defaultTSDBStatus);

    const totalParams = {
      date: requestParams.date,
      topN: 0,
      match: "",
      focusLabel: ""
    } as CardinalityRequestsParams;

    const prevDayParams = {
      ...requestParams,
      date: dayjs(requestParams.date).subtract(1, "day").tz().format(DATE_FORMAT),
    } as CardinalityRequestsParams;


    const urlBase = getCardinalityInfo(serverUrl, requestParams);
    const urlPrev = getCardinalityInfo(serverUrl, prevDayParams);
    const uslTotal = getCardinalityInfo(serverUrl, totalParams);
    const urls = [urlBase, urlPrev, uslTotal];

    try {
      const responses = await Promise.all(urls.map(url => fetch(url)));
      const [resp, respPrev, respTotals] = await Promise.all(responses.map(resp => resp.json()));
      if (responses[0].ok) {
        const { data: dataTotal } = respTotals;
        const prevResult = { ...respPrev.data } as TSDBStatus;
        const result = { ...resp.data } as TSDBStatus;
        result.totalSeriesByAll = dataTotal?.totalSeries;
        result.totalSeriesPrev = prevResult?.totalSeries;

        const name = match?.replace(/[{}"]/g, "");
        result.seriesCountByLabelValuePair = result.seriesCountByLabelValuePair.filter(s => s.name !== name);

        Object.keys(result).forEach(k => {
          const key = k as keyof TSDBStatus;
          const entries = result[key];
          const prevEntries = prevResult[key];

          if (Array.isArray(entries) && Array.isArray(prevEntries)) {
            entries.forEach((entry) => {
              const valuePrev = prevEntries.find(prevEntry => prevEntry.name === entry.name)?.value;
              entry.diff = valuePrev ? entry.value - valuePrev : 0;
              entry.valuePrev = valuePrev || 0;
            });
          }
        });

        setTSDBStatus(result);
        setIsLoading(false);
      } else {
        setError(resp.error);
        setTSDBStatus(appConfigurator.defaultTSDBStatus);
        setIsLoading(false);
      }
    } catch (e) {
      setIsLoading(false);
      if (e instanceof Error) setError(`${e.name}: ${e.message}`);
    }
  };


  useEffect(() => {
    fetchCardinalityInfo({ topN, match, date, focusLabel });
  }, [serverUrl, match, focusLabel, topN, date]);

  useEffect(() => {
    if (error) {
      setTSDBStatus(appConfigurator.defaultTSDBStatus);
      setIsLoading(false);
    }
  }, [error]);

  appConfigurator.tsdbStatusData = tsdbStatus;
  return { isLoading, appConfigurator: appConfigurator, error };
};
