import { getIronSession } from 'iron-session';
import sessionOptions from './components/sessionOptions';

export const withSession = (handler) => {
  return getIronSession(handler, sessionOptions());
};

export const withSessionApiRoute = (handler) => {
  return getIronSession(handler, sessionOptions());
};