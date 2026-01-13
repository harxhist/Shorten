import { useUI } from '../context/UIContext';


const Notification = () => {
    const { notifications, removeNotification } = useUI()
    if (!notifications || notifications.length === 0) return null;
    return(
        <div className="fixed bottom-5 right-5 flex flex-col gap-2 z-[1001]">
            {notifications.map((notification) => (
                <div key={notification.id} className={`p-3 rounded shadow-md flex items-center ${notification.type === 'error' ? 'bg-red-200 text-red-800' : 'bg-white text-black'}`}>
                    {notification.message}
                    <button onClick={() => removeNotification(notification.id)} className="ml-auto">X</button>
                </div>
            ))}
        </div>
    )
}

export default Notification;