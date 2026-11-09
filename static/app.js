document.addEventListener('DOMContentLoaded', () => {
    let currentEntity = 'emprendimientos';
    const mainForm = document.getElementById('main-form'); 
    const formFields = document.getElementById('form-fields'); 
    const formTitle = document.getElementById('form-title'); 
    const listTitle = document.getElementById('list-title'); 
    const listContainer = document.getElementById('entity-list-container'); 
    const tabButtons = document.querySelectorAll('.tab-btn');



    const formTemplates = { 
        emprendimientos: ` 
        <div class="form-group">
            <label for="nombre"> Nombre: </label>
            <input type="text" id="nombre" name="nombre" required placeholder="Ej: Panadería Tandil">
        </div>
        <div class="form-group">
            <label> RUbro </label>
            <input type="text" id="rubro" required placeholder="Ej: Gastronomía">
        </div>
        <div class="form-group">
            <label> Email </label>
            <input type="text" id="contacto" placeholder="contacto@ejemplo.com">
        </div>

        `,
        usuarios: ` 
        <div class="form-group">
            <label> Nombre Completo:</label>
            <input type="text" id="nombre_completo" required>
        </div>
        <div class="form-group">
            <label> Email </label>
            <input type="text" id="contacto" placeholder="contacto@ejemplo.com">
        </div>
        <div class="form-group">
            <label> Contraseña </label>
            <input type="password" id="password" required>
        </div>
        <div class="form-group">
            <label>Rol:</label>
            <select id="rol">
                <option value="cliente"> Cliente </option> 
                <option value="emprendedor"> Emprendedor </option>
                <option value="administrador"> Administrador </option> 
            </select>
        </div>
        `,
        publicaciones: `
            <div class="form-group">
                <label> ID emprendimiento. tiene q estar? </label>
                <input type="number" id="id_emprendimiento" required>
            </div>

            <div class="form-group">
                <label> Título </label>
                <input type="text" id="titulo" required>
            </div>

            <div class="form-group">
                <label> Tipo: </label>
                <select id="tipo" required>
                    <option value="promocion"> Promoción </option>
                    <option value="producto"> Producto </option>
                    <option value="servicio"> Servicio </option>
                </select>
            </div>
            <div class="form-group">
                <label> Precio </label>
                <input type="number" step="0.1" id="precio">
            </div>
            <div class="form-group">
                <label> Descripción  (agregar) </label>
                <textarea id="descripcion" required></textarea>
            </div>

        `,

        suscripciones: `
            <div class="form-group"> 
                <label> ID Usuario: </label>
                <input type="number" id="id_usuario" required>
            </div>
            <div class="form-group">
                <label> ID Emprendimiento: </label>
                <input type="number" id="id_emprendimiento" required>
            </div>
        `

    }       
    //cambio  de pestaña

    tabButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            tabButtons.forEach(btn => btn.classList.remove('active'));
            e.target.classList.add('active');
            currentEntity = e.target.dataset.entity;
            renderView();
        });

    });

    function renderView() {
        formFields.innerHTML = formTemplates[currentEntity];
        formTitle.textContent = `Registrar ${currentEntity.slice(0, -1)}`;
        listTitle.textContent = `Listado de ${currentEntity}`;
        fetchEntities();
    }

    async function fetchEntities() {
        listContainer.innerHTML = '<p> Cargando...</p>';
        try {
            const res = await fetch(`/api/${currentEntity}`);
            if (!res.ok) throw new Error(`Error obteniendo los datos`);
                const data = await res.json();

                if (!data || data.length === 0) {
                    listContainer.innerHTML = '<p>No hay registros guardados.</p>';
                    return;
                }

                listContainer.innerHTML = data.map(item => `
                    <div class="entity-card">
                        <div>
                            <strong>${item.nombre || item.nombre_completo || item.titulo || `Suscripcion #${item.id}`}</strong>
                            <p>{item.rubro || item.email ||item.tipo || ''}</p>   
                        </div>
                        <button class="btn-delete" data-id="${item.id || item.id_usuario}"> Eliminar</button>
                    </div>
                `).join('');
        } catch (error) {
            listContainer.innerHTML = `<p class="error-text">Error al cargar registros</p>`;
        }

    }
    

    //POST para /api/<entidad>

    mainForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const payload = {};
        const inputs = formFields.querySelectorAll('input, select');
        
        inputs.forEach(input => payload[input.id] = input.value);

        try {
            const res = await fetch(`/api/${currentEntity}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            });

            if (!res.ok || res.status === 201) {
                mainForm.reset();
                fetchEntities();
            
            } else{
                alert('Error al guardar el registro');  
            }
        } catch (err) {
            console.error('Error al enviar el formulario en POST:', err);
        }
    });

    //DELETE

    listContainer.addEventListener('click', async (e) => {
        if (e.target.classList.contains('btn-delete')) {
            const id = e.target.dataset.id;
            if (!confirm(`¿Eliminar elemento #${id}?`)) return;

            try {
                const res = await fetch(`/api/${currentEntity}/${id}`, {
                    method: 'DELETE'
                });
                
                if (res.ok || res.status === 204) {
                    fetchEntities();
                } else {
                    alert('Error al eliminar el registro');
                }

            } catch (err) {
                console.error('Error al eliminar el registro con DELETE:', err);
            }
        }
    });

renderView();

});