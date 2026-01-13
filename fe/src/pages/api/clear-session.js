import { getIronSession } from 'iron-session';
import sessionOptions from '../../components/sessionOptions';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method not allowed' });
  }

  try {
    const session = await getIronSession(req, res, sessionOptions());
    await session.destroy();
    return res.status(200).json({ message: 'Session cleared' });
  } catch (error) {
    console.error('Error clearing session:', error);
    return res.status(500).json({ message: 'Error clearing session' });
  }
}