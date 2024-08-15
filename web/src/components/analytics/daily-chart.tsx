import {
  BACKEND_URL,
  SINGLESTORE_PURPLE_700,
  SNOWFLAKE_BLUE,
} from "@/consts/config";
import { useCity, useDatabase } from "@/lib/store";
import axios from "axios";
import { useState, useEffect } from "react";
import { XAxis, YAxis, Bar, BarChart } from "recharts";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "../ui/chart";
import { Card } from "../ui/card";
import { DatabaseResultLabel } from "../ui/database-result-label";
import { format, formatDate } from "date-fns";
import { toZonedTime, fromZonedTime } from "date-fns-tz";

export default function DailyChart() {
  const database = useDatabase();
  const city = useCity();
  const [latency, setLatency] = useState(0);
  const [chartData, setChartData] = useState([]);

  useEffect(() => {
    setLatency(0);
    getData()
      .then((data) => {
        const now = new Date();
        const dailyData: { [day: string]: number } = {};
        for (let i = 0; i <= 7; i++) {
          const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000);
          const dayKey = format(date, "yyyy-MM-dd");
          dailyData[dayKey] = 0;
        }

        data.forEach((item: any) => {
          const localDate = fromZonedTime(new Date(item.daily_interval), "UTC");
          const dayKey = format(localDate, "yyyy-MM-dd");
          if (dayKey in dailyData) {
            dailyData[dayKey] = item.trip_count;
          }
        });

        const formattedData = Object.entries(dailyData).map(
          ([dayKey, trips]) => ({
            day: dayKey,
            trips: trips,
          }),
        );
        formattedData.reverse();
        setChartData(formattedData as any);
      })
      .catch((error) => console.error("Error fetching trip data:", error));
  }, [database, city]);

  const getData = async () => {
    setLatency(0);
    let cityParam = city === "All" ? "" : city;
    const response = await axios.get(
      `${BACKEND_URL}/trips/last/week?db=${database}&city=${cityParam}`,
    );
    const latencyHeader = response.headers["x-query-latency"];
    if (latencyHeader) {
      setLatency(parseInt(latencyHeader));
    }
    return response.data;
  };

  const chartConfig = {
    trips: {
      label: "Trips",
      color:
        database === "singlestore" ? SINGLESTORE_PURPLE_700 : SNOWFLAKE_BLUE,
    },
  } satisfies ChartConfig;

  return (
    <Card className="h-[400px] w-[600px]">
      <div className="flex flex-row items-center justify-between p-2">
        <h4>Ride requests per day</h4>
        <DatabaseResultLabel database={database} latency={latency} />
      </div>
      <ChartContainer config={chartConfig} className="h-full w-full pb-10 pr-4">
        <BarChart data={chartData}>
          <XAxis
            dataKey="day"
            label={{ value: "Day", position: "bottom" }}
            tickFormatter={(tick) => format(new Date(tick), "M/d")}
            interval={0}
          />
          <YAxis
            dataKey="trips"
            tickFormatter={(tick) => {
              return tick.toLocaleString();
            }}
          />
          <Bar dataKey="trips" fill="var(--color-trips)" radius={4} />
          <ChartTooltip
            content={
              <ChartTooltipContent
                labelFormatter={(value) => format(new Date(value), "M/d/yy")}
              />
            }
            cursor={false}
            defaultIndex={1}
          />
        </BarChart>
      </ChartContainer>
    </Card>
  );
}
