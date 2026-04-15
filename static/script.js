// 页面加载时获取所有待办事项
window.onload = getTodos;

// 获取所有待办事项
async function getTodos() {
    const res = await fetch('http://localhost:8080/getall');
    const todosData = await res.json();

    const ul = document.querySelector('ul');
    ul.innerHTML = '';

    if (!todosData || todosData.length === 0) {
        return;
    }

    todosData.forEach(todo => {
        ul.appendChild(createTodoItem(todo));
    });
}

// 创建待办事项的 DOM 元素
function createTodoItem(todo) {
    const li = document.createElement('li');
    li.id = `todo-${todo.id}`;
    li.innerHTML = `
        <span class="todo-text ${todo.completed ? 'completed' : ''}" id="name-${todo.id}">
            <strong class="todo-name">${todo.name}</strong>
            <span class="todo-desc">${todo.description}</span>
        </span>
        <div class="edit-area" id="edit-${todo.id}" style="display:none; flex: 1; gap: 6px; display: none;">
            <input type="text" class="edit-input" id="edit-name-${todo.id}" value="${todo.name}" placeholder="任务名称">
            <input type="text" class="edit-input" id="edit-desc-${todo.id}" value="${todo.description}" placeholder="任务描述">
        </div>
        <button class="btn-toggle" onclick="toggleTodo('${todo.id}', '${todo.name}', '${todo.description}', ${todo.completed})">
            ${todo.completed ? '撤销' : '完成'}
        </button>
        <button class="btn-edit" onclick="startEdit('${todo.id}')">更新</button>
        <button class="btn-save" id="save-${todo.id}" style="display:none;" onclick="saveEdit('${todo.id}', ${todo.completed})">保存</button>
        <button class="btn-delete" onclick="deleteTodo('${todo.id}')">删除</button>
    `;
    return li;
}

// 进入编辑模式
function startEdit(id) {
    document.getElementById(`name-${id}`).style.display = 'none';
    const editArea = document.getElementById(`edit-${id}`);
    editArea.style.display = 'flex';
    document.querySelector(`#todo-${id} .btn-edit`).style.display = 'none';
    document.getElementById(`save-${id}`).style.display = 'inline-block';
}

// 保存编辑
async function saveEdit(id, completed) {
    const name = document.getElementById(`edit-name-${id}`).value.trim();
    const desc = document.getElementById(`edit-desc-${id}`).value.trim();

    if (!name) {
        alert('任务名称不能为空');
        return;
    }

    await fetch('http://localhost:8080/update', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            id,
            name,
            description: desc,
            completed: completed.toString(),
        }),
    });

    getTodos();
}

// 创建新的待办事项
async function createTodo() {
    const name = document.getElementById('input-name').value.trim();
    const description = document.getElementById('input-description').value.trim();

    if (!name) {
        alert('任务名称不能为空');
        return;
    }

    await fetch('http://localhost:8080/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, description }),
    });

    document.getElementById('input-name').value = '';
    document.getElementById('input-description').value = '';

    getTodos();
}

// 切换完成状态
async function toggleTodo(id, name, description, completed) {
    await fetch('http://localhost:8080/update', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            id,
            name,
            description,
            completed: (!completed).toString(),
        }),
    });

    getTodos();
}

// 删除待办事项
async function deleteTodo(id) {
    await fetch('http://localhost:8080/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
    });

    getTodos();
}
