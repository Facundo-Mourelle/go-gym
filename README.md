# go-gym
Aplicación mobil para registrar mi progreso en el gimnasio un poco mas en detalle. Desarrollada utilizando Go.

## Backend

---
### Types
Analisis de Arquitectura en forma "Bottom-Up"
#### Muscle Groups
Contiene tipos para referirse a los principales grupos musculares de los que vamos a mantener registro
#### Equipment
Contiene tipos de equipamiento a utilizar para registrar de manera mas acertada los ejercicios.
Se dividen principalmente entre:
- pesos libres: el peso es directo, el perfil de resistencia depende directamente del brazo de momento
- maquinas: el peso depende del sistema empleado, ya sea poleas o maquinas de recorrido fijo, y el perfil de resistencia depende del brazo de momento, tension del sistema, y direccion de la fuerza.

Se incluye la marca del equipamiento ya que a veces a pesar de que las etiquetas de peso indiquen lo mismo, y su perfil de resistencia parezca identico, un equipamiento puede resultar mas pesado que otro. Tambien sirve como referencia rapida para distinguir equipamiento.
#### Movement Patterns
Contiene los patrones de movimiento que realizan los grupos musculares declarados, especialmente los que contribuyen principalmente a su crecimiento. Se excluyeron algunos movimientos de ciertos grupos musculares debido a su poca contribucion al crecimiento muscular.
En caso de buscar un entrenamiento mas "maximalista" en un futuro, se pueden agregar nuevos patrones de movimiento 
#### Exercise
Contiene los distintos ejercicios creados.
Para construir un ejercicio necesitamos:
- Equipamiento
- Patrones de movimiento Principales
- Patrones de movimiento Secundarios

La categorizacion en la que un patron de movimiento es "principal" o "secundario" depende de:
- Contribucion del patron de movimiento en el ejercicio
- Rango de movimiento empleado
>[!Note] Esta asignacion no es derivada automaticamente, queda a decision del usuario asignar los valores debidos al ejercicio. Esto es a fin de permitir la mayor flexibilidad a la hora de crear estos ejercicios, ya que dependen fuertemente del equipamiento utilizado

Se incluye un campo de Notas donde se pueden realizar aclaraciones adicionales sobre como realizar el ejercicio.
#### Workout
Contiene el registro de planes de entrenamiento, que agrupa ejercicios a realizar en una sesion de entrenamiento.
#### Session
Contiene el registro de las sesiones realizadas a partir de un plan de entrenamiento.
Las series de un mismo ejercicio son diferenciadas entre si para registrar mas acertadamente el progreso en cada una.
Se toman en cuenta los siguientes datos para registrar una serie:
- Repeticiones realizadas
- Repeticiones en reserva
- Numero de serie dentro de la sesion
- Peso sin procesar
- Equipamiento utilizado

A partir del peso sin procesar, se calcula el peso efectivo utilizando el perfil de resistencia asignado al equipamiento utilizado.

---
### Resistance
En este modulo nos encargamos de calcular los perfiles de resistencia asociados al equipamiento.
Los pesos libres tendran un peso efectivo directamente proporcional al peso empleado para la gran mayoria de ejercicios. En caso de añadir una variante de ejercicio que modifique el brazo de momento del ejercicio, lo introduciremos en el sector de **maquinas** (por ejemplo, un curl concentrado en banco)

Se construye el perfil de resistencia utilizando el rango del recorrido del ejercicio junto a la curva de resistencia para determinar el peso efectivo.
Las maquinas utilizan una configuracion de peso directo en vez de poleas

