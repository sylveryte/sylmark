import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { SpiderView } from "./graph/spider-view";

const queryClient = new QueryClient();
function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="h-[100vh] flex-col w-[100vw] flex">
        <SpiderView />
      </div>
    </QueryClientProvider>
  );
}

export default App;
