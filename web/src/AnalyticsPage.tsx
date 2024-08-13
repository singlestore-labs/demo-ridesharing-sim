import Header from "@/components/header";
import { Toolbar } from "@/components/toolbar";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import {
  Bar,
  BarChart,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { useEffect, useState } from "react";
import axios from "axios";
import { BACKEND_URL } from "@/consts/config";
import { useCity, useDatabase } from "@/lib/store";
import TotalStatistics from "./components/analytics/total-statistics";
import TodayStatistics from "./components/analytics/today-statistics";

const AnalyticsPage = () => {
  const database = useDatabase();
  const city = useCity();
  const [chartData, setChartData] = useState([]);

  useEffect(() => {
    let cityParam = city === "All" ? "" : city;
    axios
      .get(`${BACKEND_URL}/trips/last/day?db=${database}&city=${cityParam}`)
      .then((response) => {
        const now = new Date();
        const hourlyData: { [hour: string]: number } = {};

        // Initialize the last 24 hours with 0 trips
        for (let i = 23; i >= 0; i--) {
          const date = new Date(now.getTime() - i * 60 * 60 * 1000);
          const hourKey = date.toISOString().slice(0, 13) + ":00:00Z";
          hourlyData[hourKey] = 0;
        }

        // Fill in the data from the API response
        response.data.forEach((item: any) => {
          const hourKey =
            new Date(item.hourly_interval).toISOString().slice(0, 13) +
            ":00:00Z";
          if (hourKey in hourlyData) {
            hourlyData[hourKey] = item.trip_count;
          }
        });

        // Convert the hourlyData object to an array of objects
        const formattedData = Object.entries(hourlyData).map(
          ([hourKey, trips]) => {
            const date = new Date(hourKey);
            return {
              hour: `${date.getHours().toString().padStart(2, "0")}:00`,
              trips,
              fullDate: date.toISOString(),
            };
          },
        );

        // Sort the data by date
        formattedData.sort(
          (a, b) =>
            new Date(a.fullDate).getTime() - new Date(b.fullDate).getTime(),
        );

        setChartData(formattedData);
      })
      .catch((error) => console.error("Error fetching trip data:", error));
  }, [database, city]);

  const chartConfig = {
    trips: {
      label: "Trips",
      color: "#2563eb",
    },
  } satisfies ChartConfig;

  return (
    <div className="h-screen w-screen">
      <div className="flex w-full flex-col items-start gap-4 p-4">
        <Header currentPage="analytics" />
      </div>
      <div className="flex w-full flex-col items-start gap-4 px-4">
        <TodayStatistics />
        <TotalStatistics />
      </div>
      <ChartContainer config={chartConfig} className="min-h-[200px] w-[600px]">
        <BarChart data={chartData}>
          <XAxis
            dataKey="hour"
            label={{ value: "Hour", position: "bottom" }}
            interval={1}
          />
          <YAxis />
          <Bar dataKey="trips" fill="var(--color-trips)" radius={4} />
          <ChartTooltip
            content={<ChartTooltipContent />}
            cursor={false}
            defaultIndex={1}
          />
        </BarChart>
      </ChartContainer>
      <div className="absolute bottom-4 right-4 z-10">
        <Toolbar />
      </div>
      <div className="h-full w-full" />
    </div>
  );
};

export default AnalyticsPage;
