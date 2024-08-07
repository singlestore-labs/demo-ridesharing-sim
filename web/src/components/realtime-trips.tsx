import { useEffect, useState } from "react";
import { Card } from "@/components/ui/card";
import axios from "axios";
import { BACKEND_URL, EN_ROUTE_COLOR, SINGLESTORE_PURPLE_500, SINGLESTORE_PURPLE_700, WAITING_FOR_PICKUP_COLOR } from "@/consts/config";
import { toast } from "sonner";

interface TripStats {
  drivers_available: number;
  drivers_in_progress: number;
  riders_idle: number;
  riders_in_progress: number;
  riders_requested: number;
  riders_waiting: number;
  trips_accepted: number;
  trips_en_route: number;
  trips_requested: number;
}

interface RealtimeTripsProps {
  refreshInterval: number;
}

export function RealtimeTrips({ refreshInterval }: RealtimeTripsProps) {
  const [tripStats, setTripStats] = useState<TripStats | null>(null);

  useEffect(() => {
    const fetchTripStats = async () => {
      try {
        const response = await axios.get(`${BACKEND_URL}/trips/current`);
        setTripStats(response.data);
      } catch (error) {
        toast.error("Error refreshing trip stats");
      }
    };

    fetchTripStats();
    const interval = setInterval(fetchTripStats, refreshInterval);

    return () => clearInterval(interval);
  }, [refreshInterval]);

  if (!tripStats) return <div>...</div>;

  return (
    <div className="flex flex-wrap gap-4">
        <Card className="p-4 flex flex-col items-center justify-center">
        <h1 className="text-5xl font-bold">{tripStats.trips_requested}</h1>
        <p className="mt-2" style={{ color: SINGLESTORE_PURPLE_500 }}>Rides Requested</p>
      </Card>
      <Card className="p-4 flex flex-col items-center justify-center">
        <h1 className="text-5xl font-bold">{tripStats.drivers_available}</h1>
        <p className="mt-2" style={{ color: SINGLESTORE_PURPLE_700 }}>Drivers Available</p>
      </Card>
      <Card className="p-4 flex flex-col items-center justify-center">
        <h1 className="text-5xl font-bold">{tripStats.trips_accepted}</h1>
        <p className="mt-2" style={{ color: WAITING_FOR_PICKUP_COLOR }}>Waiting for Pickup</p>
      </Card>
      <Card className="p-4 flex flex-col items-center justify-center">
        <h1 className="text-5xl font-bold">{tripStats.riders_waiting}</h1>
        <p className="mt-2" style={{ color: WAITING_FOR_PICKUP_COLOR }}>Pickup</p>
      </Card>
      <Card className="p-4 flex flex-col items-center justify-center">
        <h1 className="text-5xl font-bold">{tripStats.riders_in_progress}</h1>
        <p className="mt-2" style={{ color: WAITING_FOR_PICKUP_COLOR }}>Dropoff</p>
      </Card>
      <Card className="p-4 flex flex-col items-center justify-center">
        <h1 className="text-5xl font-bold">{tripStats.drivers_in_progress}</h1>
        <p className="mt-2" style={{ color: EN_ROUTE_COLOR }}>In Progress</p>
      </Card>
    </div>
  );
}