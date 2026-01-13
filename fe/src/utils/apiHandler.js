import { v4 as uuidv4 } from 'uuid';
import { processContent } from './ws';
import { getCachedData, setCachedData } from './cacheUtils';

const backendUrl = process.env.NEXT_PUBLIC_API_URL;
const frontendUrl = process.env.NEXT_PUBLIC_FRONTEND_URL;

async function login(username, password) {
  const response = await fetch(`${frontendUrl}/api/auth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });

  if (response.ok) {
    const data = await response.json();
    return data.token;
  }
  console.error("[API] Login failed.");
  throw new Error('Login failed');
}

async function callApi(method, endpoint, body) {
  const token = await login('username', 'password');
  const requestId = uuidv4();
  try {
    const response = await fetch(`${backendUrl}/${endpoint}`, {
      method: method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        'X-Request-ID': requestId,
        'Origin': frontendUrl,
      },
      body: body,
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error("[API] API error:", errorText);
      throw new Error(`API error: ${errorText}`);
    }

    const data = await response.json();
    return data.responseBody;
  } catch (error) {
    console.error("[API] Error while calling API:", error);
    throw error;
  }
}

export const putFeedback = async () => {
  if (localStorage.getItem('sid')) {
    const reqID = localStorage.getItem('sid');
    let body = JSON.parse(localStorage.getItem(`feedback_${reqID}`));
    body = { ...body, requestID: reqID };

    const result = await callApi('PUT', 'feedback', JSON.stringify(body));
    return result;
  }

};

export const fetchSummaryContent = (url) => {
  return new Promise((resolve, reject) => {
    processContent(
      url,
      (content) => {
        resolve(content);
      },
      null,
      (error) => {
        console.error("[API] Error fetching summary content:", error);
        reject(error);
      }
    );
  });
};


export default callApi;