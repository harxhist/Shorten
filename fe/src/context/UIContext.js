import React, { createContext, useContext, useState } from 'react';

const UIContext = createContext();

export const UIProvider = ({ children }) => {
  const [isLoading, setIsLoading] = useState(false);
  const [notifications, setNotifications] = useState([]);

  const showLoading = () => setIsLoading(true);
  const hideLoading = () => setIsLoading(false);

  const showNotification = (message, type = 'info') => {
    const id = Date.now();
    const timeoutId = setTimeout(() => {
      setNotifications(prevNotifications => prevNotifications.filter(n => n.id !== id));
    }, 3000);
  
    setNotifications(prevNotifications => [
      ...prevNotifications,
      { id, message, type, timeoutId },
    ]);
  };
  
  const removeNotification = (id) => {
    setNotifications(prevNotifications => {
      const notification = prevNotifications.find(n => n.id === id);
      if (notification) {
        clearTimeout(notification.timeoutId);
      }
      return prevNotifications.filter(n => n.id !== id);
    });
  };
  

  const value = {
    isLoading,
    showLoading,
    hideLoading,
    notifications,
    showNotification,
    removeNotification,
  };

  return <UIContext.Provider value={value}>{children}</UIContext.Provider>;
};

export const useUI = () => {
  const context = useContext(UIContext);
  if (!context) {
    throw new Error('useUI must be used within a UIProvider');
  }
  return context;
};

const Notification = () => {
  const { notifications, removeNotification } = useUI();
  if (notifications.length === 0) return null;

  return (
    <div className="fixed bottom-5 right-5 flex flex-col gap-2 z-[1001]">
      {notifications.map((notification) => (
        <div key={notification.id} className={`p-3 rounded shadow-md flex items-center ${notification.type === 'error' ? 'bg-red-200 text-red-800' : 'bg-white text-black'}`}>
          {notification.message}
          <button onClick={() => removeNotification(notification.id)} className="ml-auto">X</button>
        </div>
      ))}
    </div>
  );
};

const Loading = () => {
  const { isLoading } = useUI();

  if (!isLoading) return null;

  return (
    <div className="fixed top-0 left-0 w-full h-full bg-black/50 backdrop-blur-md flex justify-center items-center z-[1000]">
      <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-blue-500"></div>
    </div>
  );
};

export { Notification, Loading };
