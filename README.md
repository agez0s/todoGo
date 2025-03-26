# TodoGo

TodoGo é uma API RESTful escrita em Go para gerenciar tarefas (todos) e usuários. O projeto utiliza o framework Gin para lidar com rotas HTTP, GORM para interagir com o banco de dados SQLite e JWT para autenticação. Esse é o meu primeiro projeto em Golang. Comecei com um tutorial para fazer CRUD e implementei a parte de JWT e Usuários. 
Cada Usuário possui suas Todos, e um usuário não consegue interagir com as de outro usuário. Decidi fazer desta maneira para se parecer mais com um "projeto real", visto que esse repositório foi feito apenas com fins acadêmicos.

Apesar de algumas dificuldades iniciais, consegui fazer o mínimo que eu considero "útil". Acredito que algumas partes podem ser melhor otimizadas (como a parte do JWT), e possivelmente possui algumas falhas de segurança na parte do GORM/SQLite, visto que é um ORM mas bem "cru" e em alguns casos permite um ataque de SQL Injection. 
A parte do SQL Injection é bem documentada pelo gorm, e acredito que não tenha nenhuma vulnerabilidade neste quesito.

## Funcionalidades

- **Autenticação de Usuários**:
  - Criação de novos usuários.
  - Login de usuários com geração de tokens JWT.
  - Recuperação de perfil do usuário autenticado.

- **Gerenciamento de Tarefas (Todos)**:
  - Criação de tarefas.
  - Listagem de tarefas paginada.
  - Atualização de tarefas.
  - Marcar tarefas como concluídas.
  - Exclusão de tarefas.

## Estrutura do Projeto

### Principais Diretórios

- **`config/`**: Configurações do projeto, incluindo inicialização do banco de dados e logger.
- **`db/`**: Contém o arquivo do banco de dados SQLite.
- **`docs/`**: Arquivos relacionados à documentação Swagger.
- **`handler/`**: Handlers para rotas de usuários e tarefas.
- **`router/`**: Configuração das rotas da API.
- **`schema/`**: Definição dos modelos de dados (User e Todo).
- **`utils/`**: Funções utilitárias, como geração de tokens JWT e envio de respostas HTTP.

## Pré-requisitos

- Go 1.20 ou superior.
- SQLite.
- [Postman](https://www.postman.com/) ou outra ferramenta para testar a API.

## Configuração

1. Clone o repositório:

   ```bash
   git clone https://github.com/agez0s/todoGo.git
   cd todoGo
   ```

2. Instale as dependências:

   ```bash
   go mod tidy
   ```

3. Configure as variáveis de ambiente:

   ```bash
   DBFILE=todo.db
   JWT_SECRET=seu_segredo_jwt_aqui
   ```

4. Execute o projeto:

   ```bash
   go run main.go
   ```

A API estará disponível em http://localhost:8080


## Endpoints
### Autenticação
POST **`/api/v1/auth/newUser`**: Criação de um novo usuário.\
POST **`/api/v1/auth/login`**: Login de usuário.\
GET **`/api/v1/auth/profile`**: Recupera o perfil do usuário autenticado.\

### Tarefas
POST **`/api/v1/todo/create`**: Criação de uma nova tarefa.\
GET **`/api/v1/todo/list`**: Listagem de tarefas.\
PATCH **`/api/v1/todo/update`**: Atualização de uma tarefa.\
POST **`/api/v1/todo/complete`**: Marca uma tarefa como concluída.\
DELETE **`/api/v1/todo/delete`**: Exclui uma tarefa.

## Tecnologias Utilizadas
Gin: Framework web.\
GORM: ORM para Go.\
SQLite: Banco de dados leve.\
JWT: Autenticação baseada em tokens.

## A Implementar

- Rate Limiting
- Swagger

