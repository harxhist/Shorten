import { promises as fs } from 'fs';
import path from 'path';
import jwt from 'jsonwebtoken';

const AUTH_USERNAME = process.env.AUTH_USERNAME ;
const AUTH_PASSWORD = process.env.AUTH_PASSWORD ;

const authenticateUser = (username, password) => {
  return username === AUTH_USERNAME && password === AUTH_PASSWORD;
};

let privateKey;

const loadPrivateKey = async () => {
  if (!privateKey) {
    try {
      privateKey = await fs.readFile(path.join(process.cwd(), 'keys', 'private.pem'), 'utf8');
    } catch (err) {
      throw new Error('Failed to load private key');
    }
  }
  return privateKey;
};

export default async function handler(req, res) {
  if (req.method === 'POST') {
    const { username, password } = req.body;

    if (!username || !password) {
      return res.status(400).json({ message: 'Username and password are required' });
    }
    const userIsAuthenticated = authenticateUser(username, password);

    if (userIsAuthenticated) {
      try {
        const payload = { userId: '12345', username };
        const options = { algorithm: 'RS256', expiresIn: '1m' };
        const key = await loadPrivateKey();
        const token = jwt.sign(payload, key, options);
        return res.status(200).json({ token });
      } catch (err) {
        const text = await response.text();
        console.error('Error:', text);
        return res.status(500).text({ 'message': 'Error generating token' });
      }
    } else {
      return res.status(401).json({ 'message': 'Authentication failed' });
    }
  } else {
    return res.status(405).json({ 'message': 'Method not allowed' });
  }
}
