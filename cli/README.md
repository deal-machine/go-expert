**Objetivo**: Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.


O sistema deverá gerar um relatório com informações específicas após a execução dos testes.

**Entrada de Parâmetros via CLI:**

- **--url**: URL do serviço a ser testado.
- **--requests**: Número total de requests.
- **--concurrency**: Número de chamadas simultâneas.


**Execução do Teste:**

- [x] Realizar requests HTTP para a URL especificada.
- [x] Distribuir os requests de acordo com o nível de concorrência definido.
- [x] Garantir que o número total de requests seja cumprido.

**Geração de Relatório:**

- Apresentar um relatório ao final dos testes contendo:
    - [x] Tempo total gasto na execução.
    - [x] Quantidade total de requests realizados.
    - [x] Quantidade de requests com status HTTP 200.
    - [x] Distribuição de outros códigos de status HTTP (como 404, 500, etc.).
    - [x] Relatório armazenado no arquivo logs.txt.
  
**Execução da aplicação:**
- Para rodar um container de testes `./run-server.sh`
  - Irá criar uma API na porta **3000** com endpoint */stress-test*
- Localmente
    ` go run main.go stress --url http://localhost:3000/stress-test --requests 1000 --concurrency 10`
- Build da aplicação
  - `docker build -t stress-test-cli .`
- Para rodar a aplicação via docker
  - `docker run --network=host stress-test-cli --url=http://localhost:3000/stress-test --requests=1000 --concurrency=10`
  - `docker run stress-test-cli --url=https://google.com --requests=1000 --concurrency=10`

