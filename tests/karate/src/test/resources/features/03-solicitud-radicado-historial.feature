Feature: Solicitud - radicado por la BD (C-5) e historial como única fuente de estado (C-4b)

  El radicado BNF-YYYY-NNNNNN lo genera fn_siguiente_radicado() como DEFAULT de
  la columna (no lo envía el cliente). El estado vigente de una solicitud es el
  último registro de historial_solicitud (/vigente).

  Background:
    * url crudUrl
    * def fi = hoy + 'T00:00:00Z'
    * def ff = finVigencia + 'T00:00:00Z'

  Scenario: radicado generado, historial y estado vigente
    # Beneficio prerrequisito
    Given path 'beneficio'
    And request
      """
      {
        empresa: { id: '#(empresaSeedId)' },
        categoria_beneficio_id: '#(categoriaEducacion)',
        estado_beneficio_id: '#(estadoBeneficioPublicado)',
        titulo: 'Historial C-4b - prueba Karate',
        descripcion: 'x', condiciones: 'x',
        fecha_inicio: '#(fi)', fecha_fin: '#(ff)',
        cupos_total: 5, cupos_disponibles: 5,
        usuario_creador: { id: '#(usuarioEmpresaSeedId)' }
      }
      """
    When method post
    Then status 201
    * def beneficioId = response.id

    # La solicitud se crea SIN radicado: lo pone la BD
    Given path 'solicitud-beneficio'
    And request { egresado: { id: '#(egresadoSeedId)' }, beneficio: { id: '#(beneficioId)' } }
    When method post
    Then status 201
    * def solicitudId = response.id

    Given path 'solicitud-beneficio', solicitudId
    When method get
    Then status 200
    And match response.radicado == '#regex BNF-\\d{4}-\\d{6}'

    # Sin historial todavía: /vigente responde 404 (la solicitud no tiene estado)
    Given path 'historial-solicitud/solicitud', solicitudId, 'vigente'
    When method get
    Then status 404

    # Primer registro: nace PENDIENTE (sin estado_anterior)
    Given path 'historial-solicitud'
    And request { solicitud_beneficio: { id: '#(solicitudId)' }, estado_nuevo_id: '#(estadoSolicitudPendiente)', usuario: { id: '#(usuarioEgresadoSeedId)' } }
    When method post
    Then status 201

    Given path 'historial-solicitud/solicitud', solicitudId, 'vigente'
    When method get
    Then status 200
    And match response.estado_nuevo_id == estadoSolicitudPendiente

    # Segundo registro: PENDIENTE → EN_REVISION; /vigente devuelve el más reciente
    Given path 'historial-solicitud'
    And request { solicitud_beneficio: { id: '#(solicitudId)' }, estado_anterior_id: '#(estadoSolicitudPendiente)', estado_nuevo_id: '#(estadoSolicitudEnRevision)', usuario: { id: '#(usuarioEmpresaSeedId)' }, justificacion: 'revisión iniciada por la suite' }
    When method post
    Then status 201

    Given path 'historial-solicitud/solicitud', solicitudId, 'vigente'
    When method get
    Then status 200
    And match response.estado_nuevo_id == estadoSolicitudEnRevision
    And match response.estado_anterior_id == estadoSolicitudPendiente

    # La bitácora completa llega con el más reciente primero
    Given path 'historial-solicitud/solicitud', solicitudId
    When method get
    Then status 200
    And match response == '#[2]'
    And match response[0].estado_nuevo_id == estadoSolicitudEnRevision
    And match response[1].estado_nuevo_id == estadoSolicitudPendiente
