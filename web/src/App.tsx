import { useRef, useEffect, useState, useCallback } from "react";
import mapboxgl from "mapbox-gl";
import {
  BACKEND_URL,
  EN_ROUTE_COLOR,
  MAPBOX_TOKEN,
  SINGLESTORE_PURPLE_500,
  SINGLESTORE_PURPLE_700,
  WAITING_FOR_PICKUP_COLOR,
} from "@/consts/config";
import axios from "axios";
import { toast } from "sonner";
import { Card } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ModeToggle } from "@/components/mode-toggle";
import { useTheme } from "@/components/theme-provider";
import { RealtimeTrips } from "./components/realtime-trips";
mapboxgl.accessToken = MAPBOX_TOKEN;

function App() {
  const mapContainer = useRef(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const initialLat = 50;
  const initialLong = 50;

  const [refreshInterval, setRefreshInterval] = useState(1000);
  const { theme } = useTheme();

  useEffect(() => {
    if (!mapContainer.current || map.current) return;
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style:
        theme === "dark"
          ? "mapbox://styles/mapbox/dark-v10"
          : "mapbox://styles/mapbox/light-v10",
      center: [initialLat, initialLong],
      zoom: 0,
      attributionControl: false,
    });
    map.current.on("load", () => {
      flyTo("San Francisco");
    });
  });

  useEffect(() => {
    if (map.current) {
      map.current.setStyle(
        theme === "dark"
          ? "mapbox://styles/mapbox/dark-v10"
          : "mapbox://styles/mapbox/light-v10",
      );
    }
  }, [theme]);

  const refreshData = useCallback(() => {
    const fetchData = async () => {
      try {
        await Promise.all([getRiders(), getDrivers()]);
      } catch (error) {
        toast.error("Error refreshing data");
      }
    };

    fetchData();
    const intervalId = setInterval(fetchData, refreshInterval);

    return () => clearInterval(intervalId);
  }, [refreshInterval]);

  useEffect(() => {
    const cleanup = refreshData();
    return cleanup;
  }, [refreshData]);

  const flyTo = (city: string) => {
    let coordinates = [0, 0];
    switch (city) {
      case "San Francisco":
        coordinates = [-122.4431, 37.7567];
        break;
    }
    if (map.current) {
      map.current.flyTo({
        center: [coordinates[0], coordinates[1]],
        zoom: 12,
        duration: 2000,
      });
    }
  };

  const fetchData = async (endpoint: string) => {
    try {
      const response = await axios.get(`${BACKEND_URL}/${endpoint}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching ${endpoint}:`, error);
      return [];
    }
  };

  const createGeoJSON = (
    data: any[],
    type: "riders" | "drivers",
    status: string,
  ) => ({
    type: "FeatureCollection",
    features: data
      .filter((item) => item.status === status)
      .map((item) => ({
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [item.location_long, item.location_lat],
        },
        properties: {
          id: item.id,
          name: `${item.first_name} ${item.last_name}`,
        },
      })),
  });

  const updateMapLayer = (
    map: mapboxgl.Map,
    layerId: string,
    geojson: any,
    color: string,
  ) => {
    if (map.getSource(layerId)) {
      // Update existing source
      (map.getSource(layerId) as mapboxgl.GeoJSONSource).setData(geojson);
    } else {
      // Add new source and layer
      map.addSource(layerId, {
        type: "geojson",
        data: geojson as mapboxgl.GeoJSONSourceOptions["data"],
      });

      map.addLayer({
        id: layerId,
        type: "circle",
        source: layerId,
        paint: {
          "circle-radius": 6,
          "circle-color": color,
        },
      });
    }
  };

  const getRiders = async () => {
    if (!map.current) return;

    const riders = await fetchData("riders");

    // const idleRiders = createGeoJSON(riders, 'riders', 'idle');
    // updateMapLayer(map.current, 'riders-idle', idleRiders, '#bababa');

    const requestedRiders = createGeoJSON(riders, "riders", "requested");
    updateMapLayer(
      map.current,
      "riders-requested",
      requestedRiders,
      SINGLESTORE_PURPLE_500,
    );

    const waitingRiders = createGeoJSON(riders, "riders", "waiting");
    updateMapLayer(
      map.current,
      "riders-waiting",
      waitingRiders,
      WAITING_FOR_PICKUP_COLOR,
    );
  };

  const getDrivers = async () => {
    if (!map.current) return;

    const drivers = await fetchData("drivers");

    const availableDrivers = createGeoJSON(drivers, "drivers", "available");
    updateMapLayer(
      map.current,
      "drivers-available",
      availableDrivers,
      SINGLESTORE_PURPLE_700,
    );

    const inProgressDrivers = createGeoJSON(drivers, "drivers", "in_progress");
    updateMapLayer(
      map.current,
      "drivers-in-progress",
      inProgressDrivers,
      EN_ROUTE_COLOR,
    );
  };

  return (
    <div className="relative h-screen w-screen">
      <div className="absolute left-4 top-4 z-10 flex items-center gap-2">
        <RealtimeTrips refreshInterval={refreshInterval} />
      </div>
      <div className="absolute bottom-4 right-4 z-10 flex items-center gap-2">
        <Card className="p-2">
          <div className="flex items-center gap-2">
            <p>Refresh Interval:</p>
            <Select
              onValueChange={(value) => setRefreshInterval(Number(value))}
              value={refreshInterval.toString()}
            >
              <SelectTrigger className="w-[80px]">
                <SelectValue placeholder="Refresh Interval" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1000">1s</SelectItem>
                <SelectItem value="5000">5s</SelectItem>
                <SelectItem value="10000">10s</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </Card>
        <ModeToggle />
      </div>
      <div ref={mapContainer} className="h-full w-full" />
    </div>
  );
}

export default App;
