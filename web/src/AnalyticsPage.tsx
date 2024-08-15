import Header from "@/components/header";
import { Toolbar } from "@/components/toolbar";
import TotalStatistics from "./components/analytics/total-statistics";
import TodayStatistics from "./components/analytics/today-statistics";
import HourlyChart from "./components/analytics/hourly-chart";
import DailyChart from "./components/analytics/daily-chart";

const AnalyticsPage = () => {
  return (
    <div className="h-screen w-screen">
      <div className="flex w-full flex-col items-start gap-4 p-4">
        <Header currentPage="analytics" />
      </div>
      <div className="flex w-full flex-col items-start gap-4 px-4">
        <TodayStatistics />
        <TotalStatistics />
      </div>
      <div className="flex flex-wrap items-center gap-4 p-4">
        <HourlyChart />
        <DailyChart />
      </div>
      <div className="absolute bottom-4 right-4 z-10">
        <Toolbar />
      </div>
    </div>
  );
};

export default AnalyticsPage;
