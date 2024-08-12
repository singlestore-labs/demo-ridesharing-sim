export const MAPBOX_TOKEN = import.meta.env.VITE_MAPBOX_TOKEN;
export const BACKEND_URL =
  import.meta.env.VITE_BACKEND_URL || "http://localhost:8000";

export const SINGLESTORE_PURPLE_500 = "#D199FF";
export const SINGLESTORE_PURPLE_700 = "#820DDF";
export const SINGLESTORE_PURPLE_900 = "#360061";

export const WAITING_FOR_PICKUP_COLOR = "#aaa0ad";
export const EN_ROUTE_COLOR = "#5ccc7a";

export const CITY_COORDINATES = {
  "San Francisco": {
    coordinates: [-122.4431, 37.7567],
    zoom: 12,
  },
  "San Jose": {
    coordinates: [-121.8854, 37.3382],
    zoom: 11,
  },
};
