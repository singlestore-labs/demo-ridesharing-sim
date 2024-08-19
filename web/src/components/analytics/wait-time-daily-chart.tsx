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
} from "@/components/ui/chart";
import { Card } from "@/components/ui/card";
import { DatabaseResultLabel } from "@/components/ui/database-result-label";
import { format } from "date-fns";

export default function WaitTimeDailyChart() {
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
          if (item.daily_interval in dailyData) {
            dailyData[item.daily_interval] = item.avg_wait_time;
          }
        });

        const formattedData = Object.entries(dailyData).map(
          ([dayKey, time]) => ({
            day: dayKey,
            time: time,
          }),
        );
        formattedData.reverse();
        setChartData(formattedData as any);
      })
      .catch((error) => console.error("Error fetching wait time data:", error));
  }, [database, city]);

  const getData = async () => {
    setLatency(0);
    let cityParam = city === "All" ? "" : city;
    const response = await axios.get(
      `${BACKEND_URL}/wait-time/last/week?db=${database}&city=${cityParam}`,
    );
    const latencyHeader = response.headers["x-query-latency"];
    if (latencyHeader) {
      setLatency(parseInt(latencyHeader));
    }
    return response.data;
  };

  const chartConfig = {
    time: {
      label: "Wait Time",
      color:
        database === "singlestore" ? SINGLESTORE_PURPLE_700 : SNOWFLAKE_BLUE,
    },
  } satisfies ChartConfig;

  return (
    <Card className="h-[400px] w-[600px]">
      <div className="flex flex-row items-center justify-between p-2">
        <h4>Avg rider wait time per day</h4>
        <DatabaseResultLabel database={database} latency={latency} />
      </div>
      <ChartContainer config={chartConfig} className="h-full w-full pb-10 pr-4">
        <BarChart data={chartData}>
          <XAxis
            dataKey="day"
            label={{ value: "Day", position: "bottom" }}
            tickFormatter={(tick) => {
              const [year, month, day] = tick.split("-");
              return format(
                new Date(parseInt(year), parseInt(month) - 1, parseInt(day)),
                "M/d",
              );
            }}
            interval={0}
          />
          <YAxis
            dataKey="time"
            tickFormatter={(tick) => {
              return tick.toLocaleString() + "s";
            }}
          />
          <Bar dataKey="time" fill="var(--color-time)" radius={4} />
          <ChartTooltip
            content={
              <ChartTooltipContent
                labelFormatter={(value) => {
                  const [year, month, day] = value.split("-");
                  return format(
                    new Date(
                      parseInt(year),
                      parseInt(month) - 1,
                      parseInt(day),
                    ),
                    "M/d",
                  );
                }}
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
