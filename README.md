<body>
    <h1>TestBeego Project Documentation</h1>
    <h2>Overview</h2>
    <p>The <strong>TestBeego</strong> project is a Go web application built using the Beego framework. The application integrates with external APIs and provides several functionalities such as managing cat-related data, favorites, and votes. It follows a modular structure for controllers and models.</p>
    <h2>Table of Contents</h2>
    <ul>
        <li><a href="#structure">Project Folder Structure</a></li>
        <li><a href="#setup">Project Download and Setup</a></li>
        <li><a href="#documentation">Technical Documentation</a></li>
        <li><a href="#tests">Tests</a></li>
        <li><a href="#issues">Known Issues</a></li>
    </ul>
    <h2 id="structure">Project Folder Structure</h2>
    <pre>
/home/w3e49/go/src/testBeego
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── README.md
├── testBeego/
├── conf/
│   └── app.conf
├── controllers/
│   ├── api_client_iface.go
│   ├── api_client.go
│   ├── cat_controller.go
│   ├── favorites_controller.go
│   └── vote_controller.go
├── mocks/
│   └── mock_api_client.go
├── models/
│   └── types.go
├── routers/
│   └── router.go
├── static/
│   ├── css/
│   │   └── style.css
│   ├── img/
│   │   ├── icons8-grid-view-48.png
│   │   ├── icons8-stack-48.png
│   │   └── placeholder--cat.png
│   └── js/
│       ├── main.js
│       └── reload.min.js
├── tests/
│   └── controllers_test.go
└── views/
    └── index.tpl
    </pre>
    <h2 id="setup">Project Download and Setup</h2>
    <ol>
        <li><strong>Clone the repository:</strong>
            <pre><code>git clone https://github.com/SHINO-01/testBeego.git</code></pre>
        </li>
        <li><strong>Navigate to the project directory:</strong>
            <pre><code>cd /home/user/go/src/testBeego</code></pre>
        </li>
        <li><strong>Install dependencies:</strong>
            <pre><code>go mod tidy</code></pre>
        </li>
        <li><strong>Run the application:</strong>
            <pre><code>bee run</code></pre>
        </li>
        <li>
          Access the front end at localhost:8080
        </li>
    </ol>
    <h2 id="documentation">Technical Documentation</h2>
    <p>The project uses the Beego framework, which is a Go-based framework for rapid application development.</p>
    <ul>
        <li><strong>Controllers:</strong> Manage the application's core logic and handle requests. Located in the <code>controllers/</code> directory.</li>
        <li><strong>Models:</strong> Define the data structures. Located in the <code>models/</code> directory.</li>
        <li><strong>Routers:</strong> Define the routing for HTTP requests. Located in the <code>routers/</code> directory.</li>
        <li><strong>Static:</strong> Contains static assets like CSS, JavaScript, and images.</li>
        <li><strong>Tests:</strong> Unit tests for the application are in the <code>tests/</code> directory.</li>
    </ul>
    <h2 id="tests">Tests</h2>
    <p>The project includes extensive unit tests using the <code>testify</code> framework. The test file is located in the <code>tests/</code> directory.</p>
    <ul>
        <li>Run tests with coverage:
            <pre><code>go test -cover -coverpkg=./controllers/... ./tests/...</code></pre>
        </li>
        <li>Generate a coverage report:
            <pre><code>go test -cover -coverpkg=./controllers/... ./tests/... -coverprofile=coverage.out</code></pre>
            <pre><code>go tool cover -html=coverage.out</code></pre>
        </li>
    </ul>
    <h2 id="issues">Known Issues</h2>
    <ul>
        <li>Ensure the <code>conf/app.conf</code> file exists and is properly configured. Missing configurations can lead to runtime errors.</li>
        <li>Some API endpoints may not handle malformed input correctly. more validation logic needed.</li>
    </ul>

</body>
