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
import { Skeleton } from "../ui/skeleton";
import { DatabaseResultLabel } from "../ui/database-label";

interface TripStats {
  avg_distance: number;
  avg_duration: number;
  avg_wait_time: number;
  total_trips: number;
}

export default function TotalStatistics() {
  const database = useDatabase();
  const city = useCity();

  const [tripStats, setTripStats] = useState<TripStats | null>(null);
  const [latency, setLatency] = useState(0);

  useEffect(() => {
    getTripStats();
  }, [database, city]);

  const getTripStats = async () => {
    setTripStats(null);
    setLatency(0);
    const cityParam = city === "All" ? "" : city;
    const response = await axios.get(
      `${BACKEND_URL}/trips/statistics?database=${database}&city=${cityParam}`,
    );
    setTripStats(response.data);
    const latencyHeader = response.headers["x-query-latency"];
    if (latencyHeader) {
      setLatency(parseInt(latencyHeader));
    }
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
          <h4>Lifetime Statistics</h4>
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
        <h4>Lifetime Statistics</h4>
        <DatabaseResultLabel database={database} latency={latency} />
      </div>
      <div className="mt-2 flex flex-col gap-4">
        <div className="flex flex-row flex-wrap gap-4">
          <Card className="flex flex-col items-center justify-center p-4">
            <h1 className="font-bold">
              {formatTripCount(tripStats?.total_trips)}
            </h1>
            <p
              className="mt-2 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Total Trips
            </p>
          </Card>
          <Card className="flex flex-col items-center justify-center p-4">
            <h1 className="font-bold">
              {(tripStats?.avg_distance / 1000).toFixed(3)}
            </h1>
            <p
              className="mt-2 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Avg Distance (km)
            </p>
          </Card>
          <Card className="flex flex-col items-center justify-center p-4">
            <h1 className="font-bold">
              {(tripStats?.avg_duration / 1).toFixed(1)}
            </h1>
            <p
              className="mt-2 font-medium"
              style={{ color: SINGLESTORE_PURPLE_700 }}
            >
              Avg Ride Duration (s)
            </p>
          </Card>
          <Card className="flex flex-col items-center justify-center p-4">
            <h1 className="font-bold">
              {(tripStats?.avg_wait_time / 1).toFixed(1)}
            </h1>
            <p
              className="mt-2 font-medium"
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
