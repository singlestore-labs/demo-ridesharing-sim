import createStore from 'react-superstore';

export const [useCity, setCity, getCity] = createStore("San Francisco");
export const [useDatabase, setDatabase, getDatabase] = createStore("singlestore");
export const [useRefreshInterval, setRefreshInterval, getRefreshInterval] = createStore(1000);

export const [useRiders, setRiders, getRiders] = createStore([]);
export const [useRiderLatency, setRiderLatency, getRiderLatency] = createStore(0);
export const [useDrivers, setDrivers, getDrivers] = createStore([]);
export const [useDriverLatency, setDriverLatency, getDriverLatency] = createStore(0);