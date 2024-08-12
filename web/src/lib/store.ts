import createStore from 'react-superstore';

export const [useCity, setCity, getCity] = createStore("San Francisco");
export const [useDatabase, setDatabase, getDatabase] = createStore("singlestore");
export const [useRefreshInterval, setRefreshInterval, getRefreshInterval] = createStore(1000);