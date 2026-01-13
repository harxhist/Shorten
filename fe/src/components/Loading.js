import { useUI } from '../context/UIContext';

const Loading = () => {
  const { isLoading } = useUI();

  if (!isLoading) return null;

  return (
    <div className="fixed top-0 left-0 w-full h-full bg-black/50 backdrop-blur-md flex justify-center items-center z-50">
      <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-blue-500"></div>
    </div>
  );
};

export default Loading;
