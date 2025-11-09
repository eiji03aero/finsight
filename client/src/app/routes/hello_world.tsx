import { createFileRoute } from "@tanstack/react-router";
import { HelloWorldPage } from "@/pages/HelloWorld";

export const Route = createFileRoute("/hello_world")({
	component: HelloWorldPage,
});
