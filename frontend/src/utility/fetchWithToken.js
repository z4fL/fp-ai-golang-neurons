const fetchWithToken = async (url, options = {}, token) => {
  const defaultHeaders = token ? { Authorization: `Bearer ${token}` } : {};

  const fetchOptions = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers, // Override jika ada header tambahan
    },
  };

  const response = await fetch(url, fetchOptions);
  return response;
};


export default fetchWithToken