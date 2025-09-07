import { useMutation, useQuery } from "@tanstack/react-query";
import { httpGet, httpPost } from "./api";
import type { IGraphData } from "./data";
import { Graph } from "./graph";

export const SpiderView = () => {
  const { data } = useQuery({
    queryKey: ["graph"],
    refetchOnWindowFocus: false,
    queryFn: () => httpGet<IGraphData>("graph"),
  });

  const mutate = useMutation({
    mutationFn: (id: number) =>
      httpPost("document/show", {
        id: id,
      }),
  });

  const openDoc = (id: number) => {
    mutate.mutate(id);
  };

  console.log("spider-view", data);
  return (
    <div className="overflow-hidden flex-1">
      <Graph data={data} openDoc={openDoc} />
    </div>
  );
};
