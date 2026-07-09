Feature: Beneficio - descuento y devolución ATÓMICA de cupos (RN-002b/c)

  El guard va en el propio UPDATE (cupos_disponibles > 0 al descontar;
  < cupos_total al devolver): sin condiciones de carrera y sin salirse de rango.

  Background:
    * url crudUrl
    * def fi = hoy + 'T00:00:00Z'
    * def ff = finVigencia + 'T00:00:00Z'

  Scenario: descontar hasta agotar y devolver hasta el tope
    Given path 'beneficio'
    And request
      """
      {
        empresa: { id: '#(empresaSeedId)' },
        categoria_beneficio_id: '#(categoriaEducacion)',
        estado_beneficio_id: '#(estadoBeneficioPublicado)',
        titulo: 'Cupos atómicos - prueba Karate',
        descripcion: 'beneficio para probar RN-002b/c',
        condiciones: 'ninguna',
        fecha_inicio: '#(fi)',
        fecha_fin: '#(ff)',
        cupos_total: 2,
        cupos_disponibles: 2,
        usuario_creador: { id: '#(usuarioEmpresaSeedId)' }
      }
      """
    When method post
    Then status 201
    And match response.id == '#? _ > 0'
    * def beneficioId = response.id

    # Descontar 2 veces: ambas exitosas
    Given path 'beneficio', beneficioId, 'cupo/descontar'
    And request {}
    When method post
    Then status 200
    And match response == { descontado: true }

    Given path 'beneficio', beneficioId, 'cupo/descontar'
    And request {}
    When method post
    Then status 200
    And match response == { descontado: true }

    # Tercera: sin cupos → descontado false (el guard atómico la rechaza)
    Given path 'beneficio', beneficioId, 'cupo/descontar'
    And request {}
    When method post
    Then status 200
    And match response == { descontado: false }

    Given path 'beneficio', beneficioId
    When method get
    Then status 200
    And match response.cupos_disponibles == 0

    # Devolver 2 veces: ambas exitosas; la tercera choca con el tope cupos_total
    Given path 'beneficio', beneficioId, 'cupo/devolver'
    And request {}
    When method post
    Then status 200
    And match response == { devuelto: true }

    Given path 'beneficio', beneficioId, 'cupo/devolver'
    And request {}
    When method post
    Then status 200
    And match response == { devuelto: true }

    Given path 'beneficio', beneficioId, 'cupo/devolver'
    And request {}
    When method post
    Then status 200
    And match response == { devuelto: false }

    Given path 'beneficio', beneficioId
    When method get
    Then status 200
    And match response.cupos_disponibles == 2
    And match response.cupos_total == 2

  Scenario: el borrado es lógico (activo=false), no físico
    Given path 'beneficio'
    And request
      """
      {
        empresa: { id: '#(empresaSeedId)' },
        categoria_beneficio_id: '#(categoriaEducacion)',
        estado_beneficio_id: '#(estadoBeneficioPublicado)',
        titulo: 'Borrado lógico - prueba Karate',
        descripcion: 'x', condiciones: 'x',
        fecha_inicio: '#(fi)', fecha_fin: '#(ff)',
        cupos_total: 1, cupos_disponibles: 1,
        usuario_creador: { id: '#(usuarioEmpresaSeedId)' }
      }
      """
    When method post
    Then status 201
    * def beneficioId = response.id

    Given path 'beneficio', beneficioId
    When method delete
    Then status 200

    Given path 'beneficio', beneficioId
    When method get
    Then status 200
    And match response.activo == false
