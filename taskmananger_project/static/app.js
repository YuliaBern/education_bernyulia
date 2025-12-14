const API = '/api/tasks';
const table = document.getElementById('taskTable');

async function loadTasks() {
    const res = await fetch(API);
    const tasks = await res.json();
    table.innerHTML = '';

    tasks.forEach(t => {
        const tr = document.createElement('tr');

        tr.innerHTML = `
            <td>${t.title}</td>
            <td><span class="status ${t.status}">${t.status}</span></td>
            <td class="actions">
                ${t.status === 'todo' ? `<button onclick="setStatus(${t.id}, 'in-progress')">▶</button>` : ''}
                ${t.status === 'in-progress' ? `<button onclick="setStatus(${t.id}, 'done')">✔</button>` : ''}
                <button onclick="deleteTask(${t.id})">🗑</button>
            </td>
        `;
        table.appendChild(tr);
    });
}

async function addTask() {
    const input = document.getElementById('taskInput');
    if (!input.value) return;

    await fetch(API, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ title: input.value })
    });

    input.value = '';
    loadTasks();
}

async function setStatus(id, status) {
    await fetch(`${API}/${id}/status`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ status })
    });
    loadTasks();
}

async function deleteTask(id) {
    await fetch(`${API}/${id}`, { method: 'DELETE' });
    loadTasks();
}

loadTasks();
