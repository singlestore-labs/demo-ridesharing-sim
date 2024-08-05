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
      style: "mapbox://styles/mapbox/streets-v10",
      center: [initialLat, initialLong],
      zoom: 13,
      attributionControl: false,
    });
    map.current.on('load', () => {
      getCoordinates();
    });
  });

  const getCoordinates = () => {
    if (!map.current) return;

    axios.get('http://localhost:8080/coordinates')
      .then(response => {
        const coordinates = response.data;

        // Create GeoJSON feature collection
        const geojson = {
          type: 'FeatureCollection',
          features: coordinates.map((coord: { longitude: number; latitude: number }) => ({
            type: 'Feature',
            geometry: {
              type: 'Point',
              coordinates: [coord.longitude, coord.latitude]
            },
            properties: {}
          }))
        };

        // Add GeoJSON source
        map.current!.addSource('coordinates', {
          type: 'geojson',
          data: geojson as mapboxgl.GeoJSONSourceOptions['data']
        });

        // Add layer for points
        map.current!.addLayer({
          id: 'coordinates',
          type: 'circle',
          source: 'coordinates',
          paint: {
            'circle-radius': 4,
            'circle-color': SINGLESTORE_PURPLE_700
          }
        });

        // Adjust map view to fit all points
        if (coordinates.length > 0) {
          const bounds = new mapboxgl.LngLatBounds();
          coordinates.forEach((coord: { longitude: number; latitude: number }) => {
            bounds.extend([coord.longitude, coord.latitude]);
          });
          map.current!.fitBounds(bounds, { padding: 50, duration: 500, maxZoom: 12 });
        }
      })
      .catch(error => {
        console.error('Error fetching coordinates:', error);
      });
  }

  // useEffect(() => {
  //   if (lastMessage !== null) {
  //     const data = JSON.parse(lastMessage.data);
  //     if (map.current) {
  //       if (!map.current.isMoving()) {
  //         map.current.flyTo({
  //           center: [data.longitude, data.latitude],
  //           animate: true,
  //           zoom: 16,
  //           speed: 4,
  //           bearing: data.heading,
  //         });
  //         removeMarkers();
  //         const marker = new mapboxgl.Marker(vehicleMarker())
  //           .setLngLat([data.longitude, data.latitude])
  //           .addTo(map.current);
  //         mapMarkers.push(marker);
  //       }
  //     }
  //   }
  // }, [lastMessage]);

  const vehicleMarker = () => {
    const el = document.createElement("div");
    el.className = "vehicle-marker";
    el.style.width = "40px";
    el.style.height = "40px";
    el.style.borderRadius = "50%";
    el.style.backgroundImage = "url('/icons/gps-marker-light.png')";
    el.style.backgroundSize = "cover";
    return el;
  };

  const removeMarkers = () => {
    mapMarkers.forEach((marker) => marker.remove());
    console.log("Removing markers");
    mapMarkers = [];
  };

  return (
    <div className="h-screen w-screen">
      <div ref={mapContainer} className="h-full w-full" />
    </div>
  );
}

export default App;
