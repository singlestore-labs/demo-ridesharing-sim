import { useRef, useEffect } from "react";
import mapboxgl from "mapbox-gl";
import { MAPBOX_TOKEN, SINGLESTORE_PURPLE_500, SINGLESTORE_PURPLE_700 } from "@/consts/config";
import axios from "axios";
mapboxgl.accessToken = MAPBOX_TOKEN;

function App() {
  const mapContainer = useRef(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const initialLat = 50;
  const initialLong = 50;
  let mapMarkers: mapboxgl.Marker[] = [];

  useEffect(() => {
    if (!mapContainer.current || map.current) return;
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: "mapbox://styles/mapbox/light-v10",
      center: [initialLat, initialLong],
      zoom: 0,
      attributionControl: false,
    });
    map.current.on('load', () => {
      flyTo("San Francisco");
      const refreshInterval = setInterval(() => {
        getRiders();
        getDrivers();
      }, 100);
      return () => clearInterval(refreshInterval);
    });
  });

  const flyTo = (city: string) => {
    let coordinates = [0, 0]
    switch (city) {
      case 'San Francisco':
        coordinates = [-122.4431, 37.7567]
        break;
    }
    if (map.current) {
      map.current.flyTo({ center: [coordinates[0], coordinates[1]], zoom: 12, duration: 2000 });
    }
  };

  const fetchData = async (endpoint: string) => {
    try {
      const response = await axios.get(`http://localhost:8080/${endpoint}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching ${endpoint}:`, error);
      return [];
    }
  };
  
  const createGeoJSON = (data: any[], type: 'riders' | 'drivers', status: string) => ({
    type: 'FeatureCollection',
    features: data.filter(item => item.status === status).map((item) => ({
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: [item.location.longitude, item.location.latitude]
      },
      properties: {
        id: item.id,
        name: `${item.first_name} ${item.last_name}`
      }
    }))
  });
  
  const updateMapLayer = (map: mapboxgl.Map, layerId: string, geojson: any, color: string) => {
    if (map.getSource(layerId)) {
      // Update existing source
      (map.getSource(layerId) as mapboxgl.GeoJSONSource).setData(geojson);
    } else {
      // Add new source and layer
      map.addSource(layerId, {
        type: 'geojson',
        data: geojson as mapboxgl.GeoJSONSourceOptions['data']
      });

      map.addLayer({
        id: layerId,
        type: 'circle',
        source: layerId,
        paint: {
          'circle-radius': 6,
          'circle-color': color
        }
      });
    }
  };
  
  const getRiders = async () => {
    if (!map.current) return;

    const riders = await fetchData('riders');
    
    // const idleRiders = createGeoJSON(riders, 'riders', 'idle');
    // updateMapLayer(map.current, 'riders-idle', idleRiders, '#bababa');
    
    const requestedRiders = createGeoJSON(riders, 'riders', 'requested');
    updateMapLayer(map.current, 'riders-requested', requestedRiders, SINGLESTORE_PURPLE_500);

    const waitingRiders = createGeoJSON(riders, 'riders', 'waiting');
    updateMapLayer(map.current, 'riders-waiting', waitingRiders, '#bababa');
  };
  
  const getDrivers = async () => {
    if (!map.current) return;

    const drivers = await fetchData('drivers');
    
    const availableDrivers = createGeoJSON(drivers, 'drivers', 'available');
    updateMapLayer(map.current, 'drivers-available', availableDrivers, SINGLESTORE_PURPLE_700);
    
    const inProgressDrivers = createGeoJSON(drivers, 'drivers', 'in_progress');
    updateMapLayer(map.current, 'drivers-in-progress', inProgressDrivers, '#5ed67e');
  };

  return (
    <div className="h-screen w-screen">
      <div ref={mapContainer} className="h-full w-full" />
    </div>
  );
}

export default App;