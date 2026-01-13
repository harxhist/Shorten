const CACHE_EXPIRY = 48 * 60 * 60 * 1000; // 48 hours in milliseconds

export const getCachedData = (key) => {
  try {
    const cachedItem = localStorage.getItem(key);
    if (!cachedItem) return null;

    const { data, timestamp } = JSON.parse(cachedItem);
    const now = new Date().getTime();

    if (now - timestamp > CACHE_EXPIRY) {
      localStorage.removeItem(key);
      return null;
    }

    return data;
  } catch (error) {
    console.error('[Cache] Error retrieving cached data:', error);
    return null;
  }
};

export const setCachedData = (key, data) => {
  try {
    const cacheItem = {
      data,
      timestamp: new Date().getTime()
    };
    localStorage.setItem(key, JSON.stringify(cacheItem));
  } catch (error) {
    console.error('[Cache] Error caching data:', error);
  }
};
