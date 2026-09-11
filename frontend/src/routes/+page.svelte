<script>
	let tasks = $state([]);
	let description = $state('');
	let dueDate = $state('');

	async function loadTasks() {
		const res = await fetch('http://localhost:8080/tasks');
		tasks = await res.json();
	}

	async function addTask() {
		await fetch('http://localhost:8080/tasks', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ description, due_date: dueDate })
		});

		description = '';
		dueDate = '';
		await loadTasks();
	}

	async function deleteTask(id) {
		await fetch(`http://localhost:8080/tasks/${id}`, { method: 'DELETE' });
		await loadTasks();
	}

	$effect(() => {
		loadTasks();
	});
</script>
<main>
<h1>Deadline Tracker</h1>

<form onsubmit={(e) => { e.preventDefault(); addTask(); }}>
	<input type="text" placeholder="Description" bind:value={description} required />
	<input type="date" bind:value={dueDate} required />
	<button type="submit">Add Task</button>
</form>

<ul>
	{#each tasks as task (task.id)}
		<li>
			{task.description} — due {task.due_date}
			<button onclick={() => deleteTask(task.id)}>Delete</button>
		</li>
	{/each}
</ul>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: system-ui, sans-serif;
		background: #f5f5f5;
	}

	main {
		max-width: 500px;
		margin: 0 auto;
		padding: 1rem;
	}

	h1 {
		font-size: 1.5rem;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 1.5rem;
	}

	input, button {
		padding: 0.75rem;
		font-size: 1rem;
		border-radius: 6px;
		border: 1px solid #ccc;
	}

	button {
		background: #333;
		color: white;
		border: none;
		cursor: pointer;
	}

	ul {
		list-style: none;
		padding: 0;
	}

	li {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: white;
		padding: 0.75rem;
		border-radius: 8px;
		margin-bottom: 0.5rem;
	}
</style>