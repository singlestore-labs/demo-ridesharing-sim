import Header from "@/components/header";
import TotalStatistics from "@/components/analytics/total-statistics";
import TodayStatistics from "@/components/analytics/today-statistics";
import TripsHourlyChart from "@/components/analytics/trips-hourly-chart";
import TripsDailyChart from "@/components/analytics/trips-daily-chart";
import WaitTimeDailyChart from "@/components/analytics/wait-time-daily-chart";
import WaitTimeHourlyChart from "@/components/analytics/wait-time-hourly-chart";
import WaitTimeMinuteChart from "@/components/analytics/wait-time-minute-chart";
import { Toolbar } from "@/components/toolbar";
import Pricing from "@/components/analytics/pricing";
import TripsSecondChart from "@/components/analytics/trips-second-chart";

const AnalyticsPage = () => {
  return (
    <div className="">
      <div className="flex w-full flex-col items-start gap-4 p-4">
        <Header currentPage="analytics" />
      </div>
      <div className="flex w-full flex-wrap items-start gap-4 px-4">
        <div className="flex flex-col items-start gap-4 px-4">
          <TodayStatistics />
          <TotalStatistics />
        </div>
        <div className="flex flex-col items-start gap-4 px-4">
          <Pricing />
        </div>
      </div>
      <div className="flex w-full flex-col items-start px-4">
        <div className="flex flex-col items-start p-4">
          <h4>Trends</h4>
        </div>
        <div className="flex flex-wrap items-center gap-4 px-4 pb-20">
          <TripsSecondChart />
          {/* <TripsMinuteChart /> */}
          <TripsHourlyChart />
          <TripsDailyChart />
          <WaitTimeMinuteChart />
          <WaitTimeHourlyChart />
          <WaitTimeDailyChart />
        </div>
      </div>
      <div className="fixed bottom-4 right-4 z-50">
        <Toolbar />
      </div>
    </div>
  );
};

export default AnalyticsPage;
