const LoadingSpinner = () => {
	return (
		<div className="flex h-full w-full items-center justify-center">
			<div className="rotate h-8 w-8 animate-spin rounded-full border-4 border-solid  border-l-primary" />
		</div>
	)
}

export default LoadingSpinner
