package pruebas;

import com.intuit.karate.junit5.Karate;

/**
 * Runner de la suite Karate del CRUD. No necesita mocks: el CRUD solo habla
 * con PostgreSQL. Los features se listan en orden explícito y corren en un
 * solo hilo (comparten la BD re-sembrada con db/seed_pruebas.sql).
 */
class PruebasCrudTest {

    @Karate.Test
    Karate pruebas() {
        return Karate.run(
                "classpath:features/01-getall-contrato.feature",
                "classpath:features/02-beneficio-cupos.feature",
                "classpath:features/03-solicitud-radicado-historial.feature",
                "classpath:features/04-crud-basico-usuario.feature"
        ).tags("~@ignore");
    }
}
