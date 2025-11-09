import { createRouter } from "@tanstack/react-router";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import * as TanstackQuery from "@/shared/tanstack-query";

// Import the generated route tree
import { routeTree } from "@/app/root/routeTree.gen";

// Create a new router instance
export const getRouter = () => {
	const router = createRouter({
		routeTree,
		defaultPreload: "intent",
	});

	setupRouterSsrQueryIntegration({
		router,
		queryClient: TanstackQuery.queryClient,
	});

	return router;
};
