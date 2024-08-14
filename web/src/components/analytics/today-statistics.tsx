import {
  BACKEND_URL,
  SINGLESTORE_PURPLE_500,
  SINGLESTORE_PURPLE_700,
} from "@/consts/config";
import { Card } from "@/components/ui/card";
import { useCity, useDatabase } from "@/lib/store";
import axios from "axios";
import { useState, useCallback, useEffect } from "react";
import { toast } from "sonner";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowTrendUp,
  faArrowTrendDown,
} from "@fortawesome/free-solid-svg-icons";
import { Skeleton } from "../ui/skeleton";
import { DatabaseResultLabel } from "../ui/database-label";

interface TripStats {
  total_trips_change: number;
  avg_distance_change: number;
  avg_duration_change: number;
  avg_wait_time_change: number;
  total_trips: number;
  avg_distance: number;
  avg_duration: number;
  avg_wait_time: number;
}

export default function TodayStatistics() {
  const database = useDatabase();
  const city = useCity();

  const [tripStats, setTripStats] = useState<TripStats | null>(null);
  const [latency, setLatency] = useState(0);

  useEffect(() => {
    getTripStats();
  }, [database, city]);

  const getTripStats = async () => {
    setTripStats(null);
    const cityParam = city === "All" ? "" : city;
    const response = await axios.get(
      `${BACKEND_URL}/trips/statistics/daily?database=${database}&city=${cityParam}`,
    );
    setTripStats(response.data);
  };

  const TrendDisplay = ({ change }: { change: number }) => {
    return (
      <div
        className={`flex flex-row items-center gap-2 text-sm ${
          change > 0 ? "text-green-500" : "text-red-500"
        }`}
      >
        {change > 0 ? (
          <FontAwesomeIcon icon={faArrowTrendUp} />
        ) : (
          <FontAwesomeIcon icon={faArrowTrendDown} />
        )}
        {Math.abs(change / 1).toFixed(1)}%
      </div>
    );
  };

  const formatTripCount = (count: number) => {
    if (count >= 1000000000) {
      return (count / 1000000000).toFixed(1) + "B";
    } else if (count >= 1000000) {
      return (count / 1000000).toFixed(1) + "M";
    } else if (count >= 10000) {
      return (count / 1000).toFixed(1) + "K";
    } else {
      return count.toLocaleString("en-US");
    }
  };

  if (!tripStats)
    return (
      <div>
        <div className="flex flex-row items-center justify-between">
          <h4>Today</h4>
          <DatabaseResultLabel database={database} latency={latency} />
        </div>
        <div className="mt-2 flex flex-col gap-4">
          <div className="flex flex-row flex-wrap gap-4">
            {[1, 2, 3, 4].map((_, index) => (
              <Card
                key={index}
                className="flex flex-col items-center justify-center p-4"
              >
                <Skeleton className="h-[20px] w-[100px] rounded-full" />
                <Skeleton className="mt-4 h-[20px] w-[130px] rounded-full" />
              </Card>
            ))}
          </div>
        </div>
      </div>
    );

  return (
    <div>
      <div className="flex flex-row items-center justify-between">
        <h4>Today</h4>
        <DatabaseResultLabel database={database} latency={latency} />
      </div>
      <div className="mt-2 flex flex-col gap-4">
        <div className="flex flex-row flex-wrap gap-4">
          <Card className="flex flex-col items-center justify-center py-2">
            <div className="flex w-full justify-end px-2">
              <TrendDisplay change={tripStats?.total_trips_change / 1} />
            </div>
            <h1 className="px-4 font-bold">
              {formatTripCount(tripStats?.total_trips)}
            </h1>
            <p
              className="mt-2 px-4 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Total Trips
            </p>
          </Card>
          <Card className="flex flex-col items-center justify-center py-2">
            <div className="flex w-full justify-end px-2">
              <TrendDisplay change={tripStats?.avg_distance_change / 1} />
            </div>
            <h1 className="px-4 font-bold">
              {(tripStats?.avg_distance / 1000).toFixed(3)}
            </h1>
            <p
              className="mt-2 px-4 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Avg Distance (km)
            </p>
          </Card>
          <Card className="flex flex-col items-center justify-center py-2">
            <div className="flex w-full justify-end px-2">
              <TrendDisplay change={tripStats?.avg_duration_change / 1} />
            </div>
            <h1 className="px-4 font-bold">
              {(tripStats?.avg_duration / 1).toFixed(1)}
            </h1>
            <p
              className="mt-2 px-4 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Avg Ride Duration (s)
            </p>
          </Card>
          <Card className="flex flex-col items-center justify-center py-2">
            <div className="flex w-full justify-end px-2">
              <TrendDisplay change={tripStats?.avg_wait_time_change / 1} />
            </div>
            <h1 className="px-4 font-bold">
              {(tripStats?.avg_wait_time / 1).toFixed(1)}
            </h1>
            <p
              className="mt-2 px-4 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Avg Wait Time (s)
            </p>
          </Card>
        </div>
      </div>
    </div>
  );
}
