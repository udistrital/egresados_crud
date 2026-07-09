Feature: CRUD básico de usuario (alta, consulta, actualización, borrado lógico)

  Verifica el ciclo genérico de las entidades del módulo y la restricción
  UNIQUE de documento (identidad institucional del egresado, C-7).

  Background:
    * url crudUrl
    * def documento = '9' + java.lang.System.currentTimeMillis().toString().substring(4)

  Scenario: ciclo completo de un usuario
    Given path 'usuario'
    And request { documento: '#(documento)', nombre: 'Usuario Prueba Karate', correo: 'karate@correo.udistrital.edu.co', tipo_usuario: 'EGR', sistema_origen: 'SGA', id_externo: 'KARATE-TEST-1' }
    When method post
    Then status 201
    * def usuarioId = response.id

    Given path 'usuario', usuarioId
    When method get
    Then status 200
    And match response contains { documento: '#(documento)', nombre: 'Usuario Prueba Karate', tipo_usuario: 'EGR', activo: true }

    # documento UNIQUE: repetirlo falla (violación de restricción → 500 del CRUD)
    Given path 'usuario'
    And request { documento: '#(documento)', nombre: 'Duplicado', correo: 'dup@x.co', tipo_usuario: 'EGR', sistema_origen: 'SGA' }
    When method post
    Then status 500

    # PUT reemplaza la fila completa (contrato del CRUD): se envía el objeto entero
    Given path 'usuario', usuarioId
    When method get
    Then status 200
    * def usuario = response
    * usuario.nombre = 'Nombre Actualizado Karate'
    Given path 'usuario', usuarioId
    And request usuario
    When method put
    Then status 200

    Given path 'usuario', usuarioId
    When method get
    Then status 200
    And match response.nombre == 'Nombre Actualizado Karate'

    # DELETE = borrado lógico
    Given path 'usuario', usuarioId
    When method delete
    Then status 200

    Given path 'usuario', usuarioId
    When method get
    Then status 200
    And match response.activo == false

    # Y deja de aparecer en los listados filtrados por Activo:true
    Given path 'usuario'
    And param query = 'Documento:' + documento + ',Activo:true'
    When method get
    Then status 200
    And match response == [{}]
