interface Particle {
	x: number;
	y: number;
	vx: number;
	vy: number;
	life: number;
	maxLife: number;
	color: string;
	size: number;
}

export class FireworksAnimation {
	private canvas: HTMLCanvasElement;
	private ctx: CanvasRenderingContext2D;
	private particles: Particle[] = [];
	private animationId: number | null = null;
	private resizeHandler: () => void;

	constructor(canvas: HTMLCanvasElement) {
		this.canvas = canvas;
		const ctx = canvas.getContext("2d");
		if (!ctx) throw new Error("Cannot get canvas context");
		this.ctx = ctx;

		this.resizeHandler = () => this.resizeCanvas();
		this.resizeCanvas();
		window.addEventListener("resize", this.resizeHandler);
	}

	private resizeCanvas() {
		// offsetWidthが0の場合は最小サイズを設定
		this.canvas.width = this.canvas.offsetWidth || 800;
		this.canvas.height = this.canvas.offsetHeight || 384;
	}

	private createParticles(x: number, y: number) {
		const colors = [
			"#ff0000",
			"#ff7f00",
			"#ffff00",
			"#00ff00",
			"#0000ff",
			"#4b0082",
			"#9400d3",
		];
		const particleCount = 50;

		for (let i = 0; i < particleCount; i++) {
			const angle = (Math.PI * 2 * i) / particleCount;
			const velocity = 2 + Math.random() * 3;

			this.particles.push({
				x,
				y,
				vx: Math.cos(angle) * velocity,
				vy: Math.sin(angle) * velocity,
				life: 0,
				maxLife: 60 + Math.random() * 40,
				color: colors[Math.floor(Math.random() * colors.length)],
				size: 2 + Math.random() * 2,
			});
		}
	}

	private updateParticles() {
		for (let i = this.particles.length - 1; i >= 0; i--) {
			const p = this.particles[i];

			p.x += p.vx;
			p.y += p.vy;
			p.vy += 0.1; // 重力効果
			p.vx *= 0.99; // 空気抵抗
			p.vy *= 0.99;
			p.life++;

			if (p.life > p.maxLife) {
				this.particles.splice(i, 1);
			}
		}
	}

	private drawParticles() {
		this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

		for (const p of this.particles) {
			const opacity = 1 - p.life / p.maxLife;
			this.ctx.fillStyle = p.color;
			this.ctx.globalAlpha = opacity;

			this.ctx.beginPath();
			this.ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
			this.ctx.fill();
		}

		this.ctx.globalAlpha = 1;
	}

	private animate = () => {
		this.updateParticles();
		this.drawParticles();

		if (this.particles.length > 0) {
			this.animationId = requestAnimationFrame(this.animate);
		} else {
			this.animationId = null;
		}
	};

	public fire() {
		const x = this.canvas.width / 2;
		const y = this.canvas.height / 3;

		this.createParticles(x, y);

		if (!this.animationId) {
			this.animate();
		}
	}

	public destroy() {
		if (this.animationId) {
			cancelAnimationFrame(this.animationId);
		}
		window.removeEventListener("resize", this.resizeHandler);
	}
}
