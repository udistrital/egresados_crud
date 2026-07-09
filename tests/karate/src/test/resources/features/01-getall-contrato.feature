Feature: Contrato de listado de los GetAll (scaffold de terceros_crud, regla 10)

  Todos los GetAll comparten query/fields/sortby/order/limit/offset y el idioma
  [{}] para lista vacía. Se prueba sobre /usuario (3 filas deterministas del seed).

  Background:
    * url crudUrl

  Scenario: query filtra por campo (TipoUsuario:EGR)
    Given path 'usuario'
    And param query = 'TipoUsuario:EGR,Activo:true'
    And param limit = 0
    When method get
    Then status 200
    And match response == '#[_ > 0]'
    And match each response contains { tipo_usuario: 'EGR' }
    And match response[*].documento contains '1016060113'

  Scenario: lista vacía responde el idioma [{}] del SGA
    Given path 'usuario'
    And param query = 'Documento:NO-EXISTE-999'
    When method get
    Then status 200
    And match response == [{}]

  Scenario: fields recorta columnas (llaves = nombres de campo Go)
    Given path 'usuario'
    And param query = 'Activo:true'
    And param fields = 'Id,Nombre'
    And param limit = 1
    When method get
    Then status 200
    And match response == '#[1]'
    And match response[0] == { Id: '#number', Nombre: '#string' }

  Scenario: limit y offset paginan; sortby y order ordenan
    # El seed tiene exactamente 3 usuarios (ids 1..3); las corridas crean más,
    # así que se ordena por Id ascendente y se pagina de a 2.
    Given path 'usuario'
    And param sortby = 'Id'
    And param order = 'asc'
    And param limit = 2
    When method get
    Then status 200
    And match response == '#[2]'
    And match response[0].id == 1
    And match response[1].id == 2

    Given path 'usuario'
    And param sortby = 'Id'
    And param order = 'asc'
    And param limit = 2
    And param offset = 2
    When method get
    Then status 200
    And match response[0].id == 3

    # Orden descendente: el primero ya no es el id 1
    Given path 'usuario'
    And param sortby = 'Id'
    And param order = 'desc'
    And param limit = 1
    When method get
    Then status 200
    And match response[0].id == '#? _ >= 3'

  Scenario: order sin sortby es un error de contrato
    # El scaffold heredado de terceros_crud responde los errores del listado
    # con 404 (no 400): se documenta el contrato real.
    Given path 'usuario'
    And param order = 'desc'
    When method get
    Then status 404
    And match response contains 'sortby'

  Scenario: dot-notation navega relaciones (Usuario.Id en egresado)
    Given path 'egresado'
    And param query = 'Usuario.Id:2,Activo:true'
    When method get
    Then status 200
    And match response == '#[1]'
    And match response[0].codigo_institucional == '#string'
