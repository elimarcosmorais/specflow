import type { components } from '../../shared/api/schema'
import { apiFetch } from '../../shared/api/client'

export type Project = components['schemas']['Project']
export type ProjectType = components['schemas']['ProjectType']

// O formulário envia o tipo vazio ("Selecione…") como está: quem valida é a API, e a tela
// mostra o `project.type_required` traduzido que ela devolve.
export type ProjectInput = Omit<components['schemas']['ProjectInput'], 'type'> & {
  type: ProjectType | ''
}

export async function listProjects(): Promise<Project[]> {
  const { items } = await apiFetch<{ items: Project[] }>('/api/projects')
  return items
}

export function createProject(input: ProjectInput): Promise<Project> {
  return apiFetch<Project>('/api/projects', { method: 'POST', body: input })
}

export function updateProject(id: number, input: ProjectInput): Promise<Project> {
  return apiFetch<Project>(`/api/projects/${id}`, { method: 'PUT', body: input })
}

export function deleteProject(id: number): Promise<void> {
  return apiFetch(`/api/projects/${id}`, { method: 'DELETE' })
}
