import { BACKEND_URL } from "@/consts/config";
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

export default function TestChart() {
  const database = useDatabase();
  const city = useCity();
  const [latency, setLatency] = useState(0);
  const [chartData, setChartData] = useState([]);

  useEffect(() => {
    setLatency(0);
    getData()
      .then((data) => {
        const now = new Date();
        const hourlyData: { [hour: string]: number } = {};

        // Initialize the last 24 hours with 0 trips in local time
        for (let i = 0; i <= 24; i++) {
          const date = new Date(now.getTime() - i * 60 * 60 * 1000);
          date.setMinutes(0);
          date.setSeconds(0);
          date.setMilliseconds(0);
          const hourKey = format(date, "yyyy-MM-dd HH:mm:ss");
          hourlyData[hourKey] = 0;
        }

        // Fill in the data from the API response
        data.forEach((item: any) => {
          const localDate = fromZonedTime(
            new Date(item.hourly_interval),
            "UTC",
          );
          const hourKey = format(localDate, "yyyy-MM-dd HH:mm:ss");
          if (hourKey in hourlyData) {
            hourlyData[hourKey] = item.trip_count;
          }
        });

        console.log(hourlyData);

        // Convert the hourlyData object to an array of objects
        const formattedData = Object.entries(hourlyData).map(
          ([hourKey, trips]) => ({
            hour: hourKey,
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
      `${BACKEND_URL}/trips/last/day?db=${database}&city=${cityParam}`,
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
      color: "#2563eb",
    },
  } satisfies ChartConfig;

  return (
    <Card className="h-[400px] w-[600px]">
      <div className="flex flex-row items-center justify-between p-2">
        <h4>Ride requests per hour</h4>
        <DatabaseResultLabel database={database} latency={latency} />
      </div>
      <ChartContainer config={chartConfig} className="h-full w-full">
        <BarChart data={chartData}>
          <XAxis
            dataKey="hour"
            label={{ value: "Hour", position: "bottom" }}
            tickFormatter={(tick) => format(new Date(tick), "h a")}
            interval={1}
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
                labelFormatter={(value) =>
                  format(new Date(value), "M/d h:mm a")
                }
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
