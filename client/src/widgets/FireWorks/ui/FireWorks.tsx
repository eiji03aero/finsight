"use client";

import { useEffect, useRef } from "react";
import { Button } from "@/shared/shadcn/ui/button";
import { FireworksAnimation } from "@/widgets/FireWorks/lib/fireworks";

export function FireWorks() {
	const canvasRef = useRef<HTMLCanvasElement>(null);
	const animationRef = useRef<FireworksAnimation | null>(null);

	useEffect(() => {
		if (canvasRef.current) {
			animationRef.current = new FireworksAnimation(canvasRef.current);
		}

		return () => {
			animationRef.current?.destroy();
		};
	}, []);

	const handleFire = () => {
		animationRef.current?.fire();
	};

	return (
		<div className="w-full flex flex-col gap-4">
			<canvas
				ref={canvasRef}
				className="w-full h-96 border border-border rounded-lg bg-black"
			/>
			<Button onClick={handleFire} variant="default" size="lg">
				Fire
			</Button>
		</div>
	);
}
