function fn() {
  var config = {
    crudUrl: karate.properties['crud.url'] || 'http://localhost:8080/v1',

    // Ids de los parámetros institucionales (C-1) usados por los datos de
    // prueba — mismos ids del servicio real (migración 2026-07-07):
    estadoEmpresaActiva: 7199,
    estadoBeneficioPublicado: 7202,
    estadoSolicitudPendiente: 7206,
    estadoSolicitudEnRevision: 7207,
    categoriaEducacion: 7212,

    // Ids del seed (db/seed_pruebas.sql): el script de ejecución re-siembra
    empresaSeedId: 1,
    usuarioEmpresaSeedId: 1,
    usuarioEgresadoSeedId: 2,
    egresadoSeedId: 1
  };

  var LocalDate = Java.type('java.time.LocalDate');
  config.hoy = LocalDate.now().toString();
  config.finVigencia = LocalDate.now().plusMonths(3).toString();

  karate.configure('connectTimeout', 10000);
  karate.configure('readTimeout', 30000);
  return config;
}
