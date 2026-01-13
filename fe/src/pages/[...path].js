import { getServerSidePropsForDynamicRoute } from '../utils/routeHandler';

export async function getServerSideProps(context) {
  return getServerSidePropsForDynamicRoute(context);
}

export default function DynamicPage() {
  return <p>Processing...</p>;
}