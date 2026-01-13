const sessionOptions = () => {
  const secret = process.env.SESSION_SECRET;

  if (!secret) {
    throw new Error('SESSION_SECRET environment variable is not set.');
  }

  return {
    password: secret,
    cookieName: 'Shorten',
    cookieOptions: {
      secure: process.env.NODE_ENV === 'production',
      httpOnly: true, 
      sameSite: 'strict', 
      path: '/', 
      maxAge: 86400 * 30,
    },
  };
};

export default sessionOptions;