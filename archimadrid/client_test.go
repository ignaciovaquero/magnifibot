package archimadrid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGospelFromCache(t *testing.T) {
	tests := []struct {
		name,
		key string
		cacheObject   interface{}
		expected      *Gospel
		errorExpected bool
	}{
		{
			name: "Valid gospel from Cache",
			key:  "key",
			cacheObject: &Gospel{
				Day: "today",
			},
			expected: &Gospel{
				Day: "today",
			},
			errorExpected: false,
		},
		{
			name: "Invalid object from Cache",
			key:  "key",
			cacheObject: Gospel{
				Day: "today",
			},
			errorExpected: true,
		},
		{
			name:          "No object from Cache",
			key:           "key",
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			client := NewClient()
			if test.cacheObject != nil {
				client.Set(test.key, test.cacheObject)
			}
			actual, err := client.getGospelFromCache(test.key)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.EqualValues(tt, test.expected, actual)
		})
	}
}

func TestGetResponseFromCache(t *testing.T) {
	tests := []struct {
		name,
		key string
		cacheObject   interface{}
		expected      *gospelResponse
		errorExpected bool
	}{
		{
			name: "Valid response from cache",
			key:  "key",
			cacheObject: &gospelResponse{
				PostTitle: "title",
			},
			expected: &gospelResponse{
				PostTitle: "title",
			},
			errorExpected: false,
		},
		{
			name: "Invalid response from cache",
			key:  "key",
			cacheObject: gospelResponse{
				PostTitle: "title",
			},
			errorExpected: true,
		},
		{
			name:          "No response from cache",
			key:           "key",
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			client := NewClient()
			if test.cacheObject != nil {
				client.Set(test.key, test.cacheObject)
			}
			actual, err := client.getResponseFromCache(test.key)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.EqualValues(tt, test.expected, actual)
		})
	}
}

func TestGetGospelOrLecture(t *testing.T) {
	tests := []struct {
		name string
		day  time.Time
		regexString,
		cachePrefix string
		psalm         bool
		key           string
		cache         interface{}
		response      string
		code          int
		expected      *Gospel
		errorExpected bool
	}{
		{
			name:        "Valid Gospel from cache",
			day:         time.Date(2022, time.March, 16, 0, 0, 0, 0, time.UTC),
			regexString: `(EVANGELIO).*`,
			cachePrefix: "gospel ",
			psalm:       false,
			key:         "gospel 2022-03-16",
			cache: &Gospel{
				Day: "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
			},
			response: "",
			code:     http.StatusOK,
			expected: &Gospel{
				Day: "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
			},
			errorExpected: false,
		},
		{
			name:        "Valid Gospel from cached response",
			day:         time.Date(2022, time.March, 16, 0, 0, 0, 0, time.UTC),
			regexString: `(EVANGELIO).*`,
			cachePrefix: "gospel ",
			psalm:       false,
			key:         "response 2022-03-16",
			cache: &gospelResponse{
				PostTitle:   "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
				PostContent: "<p><span class=\"Tit_Negro_Cur\">PRIMERA LECTURA</span><br /> \t\t\t\t\t\t\t<span class=\"Tit_Lectura\">Venga, vamos a hablar mal de él.</span><br /><span class=\"Tit_Negro_Normal\">Lectura del libro de Jeremías 18, 18 20</span></p>\n<p>Ellos dijeron:</p>\n<p>  «Venga, tramemos un plan contra Jeremías, porque no falta la ley del sacerdote, ni el consejo del sabio, ni el oráculo del profeta. Venga vamos a hablar mal de él y no hagamos caso de sus oráculos».</p>\n<p>Hazme caso, Señor, escucha lo que dicen mis oponentes. ¿Se paga el bien con el mal?, ¡pues me   han cavado una fosa!</p>\n<p>Recuerda que estuve ante ti, pidiendo clemencia por ellos, para apartar tu cólera.</p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.</p>\n<p></span></p>\n<p><span class=\"Tit_Lectura\">Sal 30, 5 6. 14. 15 16</span><br /><span class=\"Tit_Negro_Normal\">R. Sálvame, Señor, por tu misericordia.</span></p>\n<p>Sácame de la red que me han tendido, <br />porque tú eres mi amparo. <br />A tus manos encomiendo mi espíritu: <br />tú, el Dios leal, me librarás, R.</p>\n<p>Oigo el cuchicheo de la gente, <br />y todo me da miedo; <br />se conjuran contra mí <br />y traman quitarme la vida. R.</p>\n<p>Pero yo confío en ti, Señor, <br />te digo: «Tú eres mi Dios.» <br />En tu mano están mis azares: <br />líbrame de mis enemigos que me persiguen. R. </p>\n<p><span class=\"Tit_Lectura\">Versículo  Jn 8, 12b</span><br /><span class=\"Tit_Negro_Normal\"></span></p>\n<p>V: Yo soy la luz del mundo - dice el Señor -;<br />el que me sigue tendrá la luz de la vida. </p>\n<p><span class=\"Tit_Negro_Cur\">EVANGELIO</span><br /> \t\t\t\t\t\t<span class=\"Tit_Lectura\">Lo condenarán a muerte.</span><br /><span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio según san Mateo 20, 17-28</span></p>\n<p>En aquel tiempo, subiendo Jesús a Jerusalén, tomando aparte a los Doce, les dijo por el camino:</p>\n<p>«Mirad, estamos subiendo a Jerusalén, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenarán a muerte y lo entregarán a los gentiles, para que se burlen de él, lo azoten y lo crucifiquen; y al tercer día resucitará».</p>\n<p>Entonces se le acercó la madre de los hijos de Zebedeo con sus hijos y se postró para hacerle una petición. </p>\n<p>Él le preguntó:</p>\n<p>«¿Qué deseas?».</p>\n<p>Ella contestó:</p>\n<p>«Ordena que estos dos hijos míos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda»</p>\n<p>Pero Jesús replicó:</p>\n<p>«No sabéis lo que pedís. ¿Podéis beber el cáliz que yo he de beber?»</p>\n<p>Contestaron:</p>\n<p>«Lo somos.»</p>\n<p>Él les dijo:</p>\n<p>«Mi cáliz lo beberéis; pero sentarse a mi derecha o a mi izquierda no me toca a mí concederlo, es para aquellos para quienes lo tiene reservado mi Padre».</p>\n<p>Los otros diez, al oír aquello, se indignaron contra los dos hermanos. Y llamándolos, Jesús les dijo:</p>\n<p>«Sabéis que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No será así entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.</p>\n<p>Igual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos».</p>\n<p class=\"Tit_Negro_Normal\">Palabra del Señor.</p>\n",
			},
			response: `[{"ID":"49445","post_author":"0","post_date":"2022-03-16 00:00:00","post_date_gmt":"2022-03-15 23:00:00","post_content":"<p><span class=\"Tit_Negro_Cur\">PRIMERA LECTURA<\/span><br \/> \t\t\t\t\t\t\t<span class=\"Tit_Lectura\">Venga, vamos a hablar mal de \u00e9l.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del libro de Jerem\u00edas 18, 18 20<\/span><\/p>\n<p>Ellos dijeron:<\/p>\n<p>  \u00abVenga, tramemos un plan contra Jerem\u00edas, porque no falta la ley del sacerdote, ni el consejo del sabio, ni el or\u00e1culo del profeta. Venga vamos a hablar mal de \u00e9l y no hagamos caso de sus or\u00e1culos\u00bb.<\/p>\n<p>Hazme caso, Se\u00f1or, escucha lo que dicen mis oponentes. \u00bfSe paga el bien con el mal?, \u00a1pues me   han cavado una fosa!<\/p>\n<p>Recuerda que estuve ante ti, pidiendo clemencia por ellos, para apartar tu c\u00f3lera.<\/p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/p>\n<p><\/span><\/p>\n<p><span class=\"Tit_Lectura\">Sal 30, 5 6. 14. 15 16<\/span><br \/><span class=\"Tit_Negro_Normal\">R. S\u00e1lvame, Se\u00f1or, por tu misericordia.<\/span><\/p>\n<p>S\u00e1came de la red que me han tendido, <br \/>porque t\u00fa eres mi amparo. <br \/>A tus manos encomiendo mi esp\u00edritu: <br \/>t\u00fa, el Dios leal, me librar\u00e1s, R.<\/p>\n<p>Oigo el cuchicheo de la gente, <br \/>y todo me da miedo; <br \/>se conjuran contra m\u00ed <br \/>y traman quitarme la vida. R.<\/p>\n<p>Pero yo conf\u00edo en ti, Se\u00f1or, <br \/>te digo: \u00abT\u00fa eres mi Dios.\u00bb <br \/>En tu mano est\u00e1n mis azares: <br \/>l\u00edbrame de mis enemigos que me persiguen. R. <\/p>\n<p><span class=\"Tit_Lectura\">Vers\u00edculo  Jn 8, 12b<\/span><br \/><span class=\"Tit_Negro_Normal\"><\/span><\/p>\n<p>V: Yo soy la luz del mundo - dice el Se\u00f1or -;<br \/>el que me sigue tendr\u00e1 la luz de la vida. <\/p>\n<p><span class=\"Tit_Negro_Cur\">EVANGELIO<\/span><br \/> \t\t\t\t\t\t<span class=\"Tit_Lectura\">Lo condenar\u00e1n a muerte.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio seg\u00fan san Mateo 20, 17-28<\/span><\/p>\n<p>En aquel tiempo, subiendo Jes\u00fas a Jerusal\u00e9n, tomando aparte a los Doce, les dijo por el camino:<\/p>\n<p>\u00abMirad, estamos subiendo a Jerusal\u00e9n, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenar\u00e1n a muerte y lo entregar\u00e1n a los gentiles, para que se burlen de \u00e9l, lo azoten y lo crucifiquen; y al tercer d\u00eda resucitar\u00e1\u00bb.<\/p>\n<p>Entonces se le acerc\u00f3 la madre de los hijos de Zebedeo con sus hijos y se postr\u00f3 para hacerle una petici\u00f3n. <\/p>\n<p>\u00c9l le pregunt\u00f3:<\/p>\n<p>\u00ab\u00bfQu\u00e9 deseas?\u00bb.<\/p>\n<p>Ella contest\u00f3:<\/p>\n<p>\u00abOrdena que estos dos hijos m\u00edos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda\u00bb<\/p>\n<p>Pero Jes\u00fas replic\u00f3:<\/p>\n<p>\u00abNo sab\u00e9is lo que ped\u00eds. \u00bfPod\u00e9is beber el c\u00e1liz que yo he de beber?\u00bb<\/p>\n<p>Contestaron:<\/p>\n<p>\u00abLo somos.\u00bb<\/p>\n<p>\u00c9l les dijo:<\/p>\n<p>\u00abMi c\u00e1liz lo beber\u00e9is; pero sentarse a mi derecha o a mi izquierda no me toca a m\u00ed concederlo, es para aquellos para quienes lo tiene reservado mi Padre\u00bb.<\/p>\n<p>Los otros diez, al o\u00edr aquello, se indignaron contra los dos hermanos. Y llam\u00e1ndolos, Jes\u00fas les dijo:<\/p>\n<p>\u00abSab\u00e9is que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No ser\u00e1 as\u00ed entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.<\/p>\n<p>Igual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos\u00bb.<\/p>\n<p class=\"Tit_Negro_Normal\">Palabra del Se\u00f1or.<\/p>\n","post_title":"16\/03\/2022 - Mi\u00e9rcoles de la 2\u00aa semana de Cuaresma.","post_excerpt":"","post_status":"publish","comment_status":"open","ping_status":"open","post_password":"","post_name":"16-03-2022-miercoles-de-la-2a-semana-de-cuaresma","to_ping":"","pinged":"","post_modified":"2022-03-16 00:00:00","post_modified_gmt":"2022-03-15 23:00:00","post_content_filtered":"","post_parent":"0","guid":"https:\/\/oracionyliturgia.archimadrid.org\/?p=49445","menu_order":"0","post_type":"post","post_mime_type":"","comment_count":"3"}]`,
			code:     http.StatusOK,
			expected: &Gospel{
				Day:       "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
				Title:     "Lo condenarán a muerte.",
				Reference: "Lectura del santo Evangelio según san Mateo 20, 17-28",
				Content:   "En aquel tiempo, subiendo Jesús a Jerusalén, tomando aparte a los Doce, les dijo por el camino:\n«Mirad, estamos subiendo a Jerusalén, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenarán a muerte y lo entregarán a los gentiles, para que se burlen de él, lo azoten y lo crucifiquen; y al tercer día resucitará».\nEntonces se le acercó la madre de los hijos de Zebedeo con sus hijos y se postró para hacerle una petición. \nÉl le preguntó:\n«¿Qué deseas?».\nElla contestó:\n«Ordena que estos dos hijos míos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda»\nPero Jesús replicó:\n«No sabéis lo que pedís. ¿Podéis beber el cáliz que yo he de beber?»\nContestaron:\n«Lo somos.»\nÉl les dijo:\n«Mi cáliz lo beberéis; pero sentarse a mi derecha o a mi izquierda no me toca a mí concederlo, es para aquellos para quienes lo tiene reservado mi Padre».\nLos otros diez, al oír aquello, se indignaron contra los dos hermanos. Y llamándolos, Jesús les dijo:\n«Sabéis que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No será así entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.\nIgual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos».\n\nPalabra del Señor.",
			},
			errorExpected: false,
		},
		{
			name:        "Valid Gospel",
			day:         time.Date(2022, time.March, 16, 0, 0, 0, 0, time.UTC),
			regexString: `(EVANGELIO).*`,
			cachePrefix: "gospel ",
			psalm:       false,
			response:    `[{"ID":"49445","post_author":"0","post_date":"2022-03-16 00:00:00","post_date_gmt":"2022-03-15 23:00:00","post_content":"<p><span class=\"Tit_Negro_Cur\">PRIMERA LECTURA<\/span><br \/> \t\t\t\t\t\t\t<span class=\"Tit_Lectura\">Venga, vamos a hablar mal de \u00e9l.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del libro de Jerem\u00edas 18, 18 20<\/span><\/p>\n<p>Ellos dijeron:<\/p>\n<p>  \u00abVenga, tramemos un plan contra Jerem\u00edas, porque no falta la ley del sacerdote, ni el consejo del sabio, ni el or\u00e1culo del profeta. Venga vamos a hablar mal de \u00e9l y no hagamos caso de sus or\u00e1culos\u00bb.<\/p>\n<p>Hazme caso, Se\u00f1or, escucha lo que dicen mis oponentes. \u00bfSe paga el bien con el mal?, \u00a1pues me   han cavado una fosa!<\/p>\n<p>Recuerda que estuve ante ti, pidiendo clemencia por ellos, para apartar tu c\u00f3lera.<\/p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/p>\n<p><\/span><\/p>\n<p><span class=\"Tit_Lectura\">Sal 30, 5 6. 14. 15 16<\/span><br \/><span class=\"Tit_Negro_Normal\">R. S\u00e1lvame, Se\u00f1or, por tu misericordia.<\/span><\/p>\n<p>S\u00e1came de la red que me han tendido, <br \/>porque t\u00fa eres mi amparo. <br \/>A tus manos encomiendo mi esp\u00edritu: <br \/>t\u00fa, el Dios leal, me librar\u00e1s, R.<\/p>\n<p>Oigo el cuchicheo de la gente, <br \/>y todo me da miedo; <br \/>se conjuran contra m\u00ed <br \/>y traman quitarme la vida. R.<\/p>\n<p>Pero yo conf\u00edo en ti, Se\u00f1or, <br \/>te digo: \u00abT\u00fa eres mi Dios.\u00bb <br \/>En tu mano est\u00e1n mis azares: <br \/>l\u00edbrame de mis enemigos que me persiguen. R. <\/p>\n<p><span class=\"Tit_Lectura\">Vers\u00edculo  Jn 8, 12b<\/span><br \/><span class=\"Tit_Negro_Normal\"><\/span><\/p>\n<p>V: Yo soy la luz del mundo - dice el Se\u00f1or -;<br \/>el que me sigue tendr\u00e1 la luz de la vida. <\/p>\n<p><span class=\"Tit_Negro_Cur\">EVANGELIO<\/span><br \/> \t\t\t\t\t\t<span class=\"Tit_Lectura\">Lo condenar\u00e1n a muerte.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio seg\u00fan san Mateo 20, 17-28<\/span><\/p>\n<p>En aquel tiempo, subiendo Jes\u00fas a Jerusal\u00e9n, tomando aparte a los Doce, les dijo por el camino:<\/p>\n<p>\u00abMirad, estamos subiendo a Jerusal\u00e9n, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenar\u00e1n a muerte y lo entregar\u00e1n a los gentiles, para que se burlen de \u00e9l, lo azoten y lo crucifiquen; y al tercer d\u00eda resucitar\u00e1\u00bb.<\/p>\n<p>Entonces se le acerc\u00f3 la madre de los hijos de Zebedeo con sus hijos y se postr\u00f3 para hacerle una petici\u00f3n. <\/p>\n<p>\u00c9l le pregunt\u00f3:<\/p>\n<p>\u00ab\u00bfQu\u00e9 deseas?\u00bb.<\/p>\n<p>Ella contest\u00f3:<\/p>\n<p>\u00abOrdena que estos dos hijos m\u00edos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda\u00bb<\/p>\n<p>Pero Jes\u00fas replic\u00f3:<\/p>\n<p>\u00abNo sab\u00e9is lo que ped\u00eds. \u00bfPod\u00e9is beber el c\u00e1liz que yo he de beber?\u00bb<\/p>\n<p>Contestaron:<\/p>\n<p>\u00abLo somos.\u00bb<\/p>\n<p>\u00c9l les dijo:<\/p>\n<p>\u00abMi c\u00e1liz lo beber\u00e9is; pero sentarse a mi derecha o a mi izquierda no me toca a m\u00ed concederlo, es para aquellos para quienes lo tiene reservado mi Padre\u00bb.<\/p>\n<p>Los otros diez, al o\u00edr aquello, se indignaron contra los dos hermanos. Y llam\u00e1ndolos, Jes\u00fas les dijo:<\/p>\n<p>\u00abSab\u00e9is que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No ser\u00e1 as\u00ed entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.<\/p>\n<p>Igual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos\u00bb.<\/p>\n<p class=\"Tit_Negro_Normal\">Palabra del Se\u00f1or.<\/p>\n","post_title":"16\/03\/2022 - Mi\u00e9rcoles de la 2\u00aa semana de Cuaresma.","post_excerpt":"","post_status":"publish","comment_status":"open","ping_status":"open","post_password":"","post_name":"16-03-2022-miercoles-de-la-2a-semana-de-cuaresma","to_ping":"","pinged":"","post_modified":"2022-03-16 00:00:00","post_modified_gmt":"2022-03-15 23:00:00","post_content_filtered":"","post_parent":"0","guid":"https:\/\/oracionyliturgia.archimadrid.org\/?p=49445","menu_order":"0","post_type":"post","post_mime_type":"","comment_count":"3"}]`,
			code:        http.StatusOK,
			expected: &Gospel{
				Day:       "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
				Title:     "Lo condenarán a muerte.",
				Reference: "Lectura del santo Evangelio según san Mateo 20, 17-28",
				Content:   "En aquel tiempo, subiendo Jesús a Jerusalén, tomando aparte a los Doce, les dijo por el camino:\n«Mirad, estamos subiendo a Jerusalén, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenarán a muerte y lo entregarán a los gentiles, para que se burlen de él, lo azoten y lo crucifiquen; y al tercer día resucitará».\nEntonces se le acercó la madre de los hijos de Zebedeo con sus hijos y se postró para hacerle una petición. \nÉl le preguntó:\n«¿Qué deseas?».\nElla contestó:\n«Ordena que estos dos hijos míos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda»\nPero Jesús replicó:\n«No sabéis lo que pedís. ¿Podéis beber el cáliz que yo he de beber?»\nContestaron:\n«Lo somos.»\nÉl les dijo:\n«Mi cáliz lo beberéis; pero sentarse a mi derecha o a mi izquierda no me toca a mí concederlo, es para aquellos para quienes lo tiene reservado mi Padre».\nLos otros diez, al oír aquello, se indignaron contra los dos hermanos. Y llamándolos, Jesús les dijo:\n«Sabéis que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No será así entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.\nIgual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos».\n\nPalabra del Señor.",
			},
			errorExpected: false,
		},
		{
			name:        "Valid First Lecture",
			day:         time.Date(2022, time.March, 16, 0, 0, 0, 0, time.UTC),
			regexString: `(PRIMERA\sLECTURA).*?Palabra de Dios\.`,
			cachePrefix: "first lecture ",
			psalm:       false,
			response:    `[{"ID":"49445","post_author":"0","post_date":"2022-03-16 00:00:00","post_date_gmt":"2022-03-15 23:00:00","post_content":"<p><span class=\"Tit_Negro_Cur\">PRIMERA LECTURA<\/span><br \/> \t\t\t\t\t\t\t<span class=\"Tit_Lectura\">Venga, vamos a hablar mal de \u00e9l.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del libro de Jerem\u00edas 18, 18 20<\/span><\/p>\n<p>Ellos dijeron:<\/p>\n<p>  \u00abVenga, tramemos un plan contra Jerem\u00edas, porque no falta la ley del sacerdote, ni el consejo del sabio, ni el or\u00e1culo del profeta. Venga vamos a hablar mal de \u00e9l y no hagamos caso de sus or\u00e1culos\u00bb.<\/p>\n<p>Hazme caso, Se\u00f1or, escucha lo que dicen mis oponentes. \u00bfSe paga el bien con el mal?, \u00a1pues me   han cavado una fosa!<\/p>\n<p>Recuerda que estuve ante ti, pidiendo clemencia por ellos, para apartar tu c\u00f3lera.<\/p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/p>\n<p><\/span><\/p>\n<p><span class=\"Tit_Lectura\">Sal 30, 5 6. 14. 15 16<\/span><br \/><span class=\"Tit_Negro_Normal\">R. S\u00e1lvame, Se\u00f1or, por tu misericordia.<\/span><\/p>\n<p>S\u00e1came de la red que me han tendido, <br \/>porque t\u00fa eres mi amparo. <br \/>A tus manos encomiendo mi esp\u00edritu: <br \/>t\u00fa, el Dios leal, me librar\u00e1s, R.<\/p>\n<p>Oigo el cuchicheo de la gente, <br \/>y todo me da miedo; <br \/>se conjuran contra m\u00ed <br \/>y traman quitarme la vida. R.<\/p>\n<p>Pero yo conf\u00edo en ti, Se\u00f1or, <br \/>te digo: \u00abT\u00fa eres mi Dios.\u00bb <br \/>En tu mano est\u00e1n mis azares: <br \/>l\u00edbrame de mis enemigos que me persiguen. R. <\/p>\n<p><span class=\"Tit_Lectura\">Vers\u00edculo  Jn 8, 12b<\/span><br \/><span class=\"Tit_Negro_Normal\"><\/span><\/p>\n<p>V: Yo soy la luz del mundo - dice el Se\u00f1or -;<br \/>el que me sigue tendr\u00e1 la luz de la vida. <\/p>\n<p><span class=\"Tit_Negro_Cur\">EVANGELIO<\/span><br \/> \t\t\t\t\t\t<span class=\"Tit_Lectura\">Lo condenar\u00e1n a muerte.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio seg\u00fan san Mateo 20, 17-28<\/span><\/p>\n<p>En aquel tiempo, subiendo Jes\u00fas a Jerusal\u00e9n, tomando aparte a los Doce, les dijo por el camino:<\/p>\n<p>\u00abMirad, estamos subiendo a Jerusal\u00e9n, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenar\u00e1n a muerte y lo entregar\u00e1n a los gentiles, para que se burlen de \u00e9l, lo azoten y lo crucifiquen; y al tercer d\u00eda resucitar\u00e1\u00bb.<\/p>\n<p>Entonces se le acerc\u00f3 la madre de los hijos de Zebedeo con sus hijos y se postr\u00f3 para hacerle una petici\u00f3n. <\/p>\n<p>\u00c9l le pregunt\u00f3:<\/p>\n<p>\u00ab\u00bfQu\u00e9 deseas?\u00bb.<\/p>\n<p>Ella contest\u00f3:<\/p>\n<p>\u00abOrdena que estos dos hijos m\u00edos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda\u00bb<\/p>\n<p>Pero Jes\u00fas replic\u00f3:<\/p>\n<p>\u00abNo sab\u00e9is lo que ped\u00eds. \u00bfPod\u00e9is beber el c\u00e1liz que yo he de beber?\u00bb<\/p>\n<p>Contestaron:<\/p>\n<p>\u00abLo somos.\u00bb<\/p>\n<p>\u00c9l les dijo:<\/p>\n<p>\u00abMi c\u00e1liz lo beber\u00e9is; pero sentarse a mi derecha o a mi izquierda no me toca a m\u00ed concederlo, es para aquellos para quienes lo tiene reservado mi Padre\u00bb.<\/p>\n<p>Los otros diez, al o\u00edr aquello, se indignaron contra los dos hermanos. Y llam\u00e1ndolos, Jes\u00fas les dijo:<\/p>\n<p>\u00abSab\u00e9is que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No ser\u00e1 as\u00ed entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.<\/p>\n<p>Igual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos\u00bb.<\/p>\n<p class=\"Tit_Negro_Normal\">Palabra del Se\u00f1or.<\/p>\n","post_title":"16\/03\/2022 - Mi\u00e9rcoles de la 2\u00aa semana de Cuaresma.","post_excerpt":"","post_status":"publish","comment_status":"open","ping_status":"open","post_password":"","post_name":"16-03-2022-miercoles-de-la-2a-semana-de-cuaresma","to_ping":"","pinged":"","post_modified":"2022-03-16 00:00:00","post_modified_gmt":"2022-03-15 23:00:00","post_content_filtered":"","post_parent":"0","guid":"https:\/\/oracionyliturgia.archimadrid.org\/?p=49445","menu_order":"0","post_type":"post","post_mime_type":"","comment_count":"3"}]`,
			code:        http.StatusOK,
			expected: &Gospel{
				Day:       "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
				Title:     "Venga, vamos a hablar mal de él.",
				Reference: "Lectura del libro de Jeremías 18, 18 20",
				Content:   "Ellos dijeron:\n  «Venga, tramemos un plan contra Jeremías, porque no falta la ley del sacerdote, ni el consejo del sabio, ni el oráculo del profeta. Venga vamos a hablar mal de él y no hagamos caso de sus oráculos».\nHazme caso, Señor, escucha lo que dicen mis oponentes. ¿Se paga el bien con el mal?, ¡pues me   han cavado una fosa!\nRecuerda que estuve ante ti, pidiendo clemencia por ellos, para apartar tu cólera.\n\nPalabra de Dios.",
			},
			errorExpected: false,
		},
		{
			name:        "Valid Second Lecture",
			day:         time.Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC),
			regexString: `(SEGUNDA\sLECTURA).*?Palabra de Dios\.`,
			cachePrefix: "second lecture ",
			psalm:       false,
			response:    `[{"ID":"49449","post_author":"0","post_date":"2022-03-20 00:00:00","post_date_gmt":"2022-03-19 23:00:00","post_content":"<p><span class=\"Tit_Negro_Cur\">PRIMERA LECTURA<\/span><br \/> \t\t\t\t\t\t\t<span class=\"Tit_Lectura\">\u201cYo soy\u201d me env\u00eda a vosotros.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del libro del \u00c9xodo 3, 1-8a. 13-15<\/span><\/p>\n<p>En aquellos d\u00edas, Mois\u00e9s pastoreaba el reba\u00f1o de su suegro Jetr\u00f3, sacerdote de Madi\u00e1n. Llev\u00f3 el reba\u00f1o trashumando por el desierto hasta llegar a Horeb, la monta\u00f1a de Dios.<\/p>\n<p>El \u00e1ngel del Se\u00f1or se le apareci\u00f3 en una llamarada entre las zarzas. Mois\u00e9s se fij\u00f3: la zarza ard\u00eda sin consumirse.<\/p>\n<p>Mois\u00e9s se dijo:<\/p>\n<p>\u00abVoy a acercarme a mirar este espect\u00e1culo admirable, a ver por qu\u00e9 no se quema la zarza\u00bb.<\/p>\n<p>Viendo el Se\u00f1or que Mois\u00e9s se acercaba a mirar, lo llam\u00f3 desde la zarza:<\/p>\n<p>\u00abMois\u00e9s, Mois\u00e9s\u00bb<\/p>\n<p>Respondi\u00f3 \u00e9l:<\/p>\n<p>\u00abAqu\u00ed estoy\u00bb<\/p>\n<p>Dijo Dios:<\/p>\n<p>\u00abNo te acerques; qu\u00edtate las sandalias de los pies, pues el sitio que pisas es terreno sagrado\u00bb.<\/p>\n<p>Y a\u00f1adi\u00f3:<\/p>\n<p>\u00abYo soy el Dios de tus padres, el Dios de Abrah\u00e1n, el Dios de Isaac, el Dios de Jacob\u00bb<\/p>\n<p>Mois\u00e9s se tap\u00f3 la cara, porque  tem\u00eda ver a Dios.<\/p>\n<p>El Se\u00f1or le dijo:<\/p>\n<p>\u00abHe visto la opresi\u00f3n de mi pueblo en Egipto y he o\u00eddo sus quejas contra los opresores, conozco sus sufrimientos. He bajado a librarlo de los egipcios, a sacarlo de esta tierra, para llevarlo a una tierra f\u00e9rtil y espaciosa, tierra que mana leche y miel\u00bb<\/p>\n<p>Mois\u00e9s replic\u00f3 a Dios:<\/p>\n<p>\u00abMira, yo ir\u00e9 a los hijos de Israel y les dir\u00e9: \u201cEl Dios de vuestros padres me ha enviado a vosotros\u201d. Si ellos me preguntan: \u201c\u00bfCu\u00e1l es su nombre? \u201c, \u00bfqu\u00e9 les respondo?\u00bb<\/p>\n<p>Dios dijo a Mois\u00e9s:<\/p>\n<p>\u00ab\u201cYo soy el que soy\u201d; esto dir\u00e1s a los hijos de Israel: \u201cYo soy\u201d me env\u00eda a vosotros\u00bb.<\/p>\n<p>Dios a\u00f1adi\u00f3:<\/p>\n<p>\u00abEsto dir\u00e1s a los hijos de Israel: \u201cEl Se\u00f1or, Dios de vuestros padres, el Dios de Abrah\u00e1n, Dios de Isaac, Dios de Jacob, me env\u00eda a vosotros. Este es mi nombre para siempre: as\u00ed me llamar\u00e9is de generaci\u00f3n en generaci\u00f3n\u201d\u00bb.<\/p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/p>\n<p><\/span><\/p>\n<p><span class=\"Tit_Lectura\">Sal 102, 1-2. 3-4. 6-7. 8 y 11<\/span><br \/><span class=\"Tit_Negro_Normal\">R. El Se\u00f1or es compasivo y misericordioso.<\/span><\/p>\n<p>Bendice, alma m\u00eda, al Se\u00f1or, <br \/>y todo mi ser a su santo nombre. <br \/>Bendice, alma m\u00eda, al Se\u00f1or, <br \/>y no olvides sus beneficios. R.<\/p>\n<p>\u00c9l perdona todas tus culpas <br \/>y cura todas tus enfermedades; <br \/>\u00e9l rescata tu vida de la fosa <br \/>y te colma de gracia y de ternura. R.<\/p>\n<p>El Se\u00f1or hace justicia <br \/>y defiende a todos los oprimidos; <br \/>ense\u00f1\u00f3 sus caminos a Mois\u00e9s <br \/>y sus haza\u00f1as a los hijos de Israel. R.<\/p>\n<p>El Se\u00f1or es compasivo y misericordioso,<br \/>lento a la ira y rico en clemencia. <br \/>Como se levanta el cielo sobre la tierra, <br \/>se levanta su bondad sobre los que lo temen. R. <\/p>\n<p><span class=\"Tit_Negro_Cur\">SEGUNDA LECTURA<\/span><br \/><span class=\"Tit_Lectura\">La vida del pueblo con Mois\u00e9s en el desierto fue escrita para escarmiento nuestro.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura de la primera carta del ap\u00f3stol san Pablo a los Corintios 10, 1-6. 10-12<\/span><\/p>\n<p>No quiero que ignor\u00e9is, hermanos, que nuestros padres estuvieron todos bajo la nube y todos atravesaron el mar y todos fueron bautizados en Mois\u00e9s por la nube y por el mar y todos comieron el mismo alimento espiritual; y todos bebieron la misma bebida espiritual, pues beb\u00edan de la roca espiritual que los segu\u00eda; y la roca era Cristo. Pero la mayor\u00eda de ellos no agradaron a Dios, pues sus cuerpos quedaron tendidos en el desierto.<\/p>\n<p>Estas cosas sucedieron en figura para nosotros, para que no codiciemos el mal como lo codiciaron ellos. Y para que no murmur\u00e9is. como murmuraron algunos de ellos, y  perecieron a manos del Exterminador.<\/p>\n<p>Todo esto les suced\u00eda aleg\u00f3ricamente y fue escrito para escarmiento nuestro, a quienes nos ha tocado vivir en la \u00faltima de las edades. Por lo tanto, el que se crea seguro, cu\u00eddese de no caer.<\/p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/span><\/p>\n<p><span class=\"Tit_Lectura\">Vers\u00edculo  Mt 4, 17<\/span><br \/><span class=\"Tit_Negro_Normal\"><\/span><\/p>\n<p>V: Convert\u00edos - dice el se\u00f1or -, <br \/>porque est\u00e1 cerca el reino de los cielos. <\/p>\n<p><span class=\"Tit_Negro_Cur\">EVANGELIO<\/span><br \/> \t\t\t\t\t\t<span class=\"Tit_Lectura\">Si no os convert\u00eds, todos perecer\u00e9is de la misma manera.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio seg\u00fan san Lucas 13, 1-9<\/span><\/p>\n<p>En aquel momento se presentaron algunos a contar a Jes\u00fas lo de los galileos, cuya sangre hab\u00eda mezclado Pilato con la de los sacrificios que ofrec\u00edan. <\/p>\n<p>Jes\u00fas respondi\u00f3:<\/p>\n<p> \u00ab \u00bfPens\u00e1is que esos galileos eran m\u00e1s pecadores que los dem\u00e1s galileos porque han padecido todo esto? Os digo que no; y, si no os convert\u00eds, todos perecer\u00e9is lo mismo. O aquellos dieciocho sobre los que cay\u00f3  la torre de Silo\u00e9 y los mat\u00f3, \u00bfpens\u00e1is que eran m\u00e1s culpables que los dem\u00e1s habitantes de Jerusal\u00e9n? Os digo que no; y, si no os convert\u00eds, todos perecer\u00e9is de la misma manera\u00bb.<\/p>\n<p>Y les dijo esta par\u00e1bola:<\/p>\n<p>\u00abUno ten\u00eda una higuera plantada en su vi\u00f1a, y fue a buscar fruto en ella, y no lo encontr\u00f3.<\/p>\n<p>Dijo entonces al vi\u00f1ador:<\/p>\n<p>\"Ya ves, tres a\u00f1os llevo viniendo a buscar fruto en esta higuera, y no lo encuentro. C\u00f3rtala. \u00bfPara qu\u00e9 va a perjudicar el terreno?\".<\/p>\n<p>Pero el vi\u00f1ador contest\u00f3:<\/p>\n<p>\"Se\u00f1or, d\u00e9jala todav\u00eda este a\u00f1o y mientras tanto yo cavar\u00e9 alrededor y le echar\u00e9 esti\u00e9rcol, a ver si da fruto en adelante. Si no, la puedes cortar\"\u00bb.<\/p>\n<p class=\"Tit_Negro_Normal\">Palabra del Se\u00f1or.<\/p>\n","post_title":"20\/03\/2022 - Domingo de la 3\u00aa semana de Cuaresma.","post_excerpt":"","post_status":"future","comment_status":"open","ping_status":"open","post_password":"","post_name":"20-03-2022-domingo-de-la-3a-semana-de-cuaresma","to_ping":"","pinged":"","post_modified":"2022-03-20 00:00:00","post_modified_gmt":"2022-03-19 23:00:00","post_content_filtered":"","post_parent":"0","guid":"https:\/\/oracionyliturgia.archimadrid.org\/?p=49449","menu_order":"0","post_type":"post","post_mime_type":"","comment_count":"0"}]`,
			code:        http.StatusOK,
			expected: &Gospel{
				Day:       "20/03/2022 - Domingo de la 3ª semana de Cuaresma.",
				Title:     "La vida del pueblo con Moisés en el desierto fue escrita para escarmiento nuestro.",
				Reference: "Lectura de la primera carta del apóstol san Pablo a los Corintios 10, 1-6. 10-12",
				Content:   "No quiero que ignoréis, hermanos, que nuestros padres estuvieron todos bajo la nube y todos atravesaron el mar y todos fueron bautizados en Moisés por la nube y por el mar y todos comieron el mismo alimento espiritual; y todos bebieron la misma bebida espiritual, pues bebían de la roca espiritual que los seguía; y la roca era Cristo. Pero la mayoría de ellos no agradaron a Dios, pues sus cuerpos quedaron tendidos en el desierto.\nEstas cosas sucedieron en figura para nosotros, para que no codiciemos el mal como lo codiciaron ellos. Y para que no murmuréis. como murmuraron algunos de ellos, y  perecieron a manos del Exterminador.\nTodo esto les sucedía alegóricamente y fue escrito para escarmiento nuestro, a quienes nos ha tocado vivir en la última de las edades. Por lo tanto, el que se crea seguro, cuídese de no caer.\n\nPalabra de Dios.",
			},
			errorExpected: false,
		},
		{
			name:        "Valid Psalm",
			day:         time.Date(2022, time.March, 16, 0, 0, 0, 0, time.UTC),
			regexString: `(Palabra\sde\sDios\..*<p>)<span.*?\sR.\s`,
			cachePrefix: "psalm ",
			psalm:       true,
			response:    `[{"ID":"49445","post_author":"0","post_date":"2022-03-16 00:00:00","post_date_gmt":"2022-03-15 23:00:00","post_content":"<p><span class=\"Tit_Negro_Cur\">PRIMERA LECTURA<\/span><br \/> \t\t\t\t\t\t\t<span class=\"Tit_Lectura\">Venga, vamos a hablar mal de \u00e9l.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del libro de Jerem\u00edas 18, 18 20<\/span><\/p>\n<p>Ellos dijeron:<\/p>\n<p>  \u00abVenga, tramemos un plan contra Jerem\u00edas, porque no falta la ley del sacerdote, ni el consejo del sabio, ni el or\u00e1culo del profeta. Venga vamos a hablar mal de \u00e9l y no hagamos caso de sus or\u00e1culos\u00bb.<\/p>\n<p>Hazme caso, Se\u00f1or, escucha lo que dicen mis oponentes. \u00bfSe paga el bien con el mal?, \u00a1pues me   han cavado una fosa!<\/p>\n<p>Recuerda que estuve ante ti, pidiendo clemencia por ellos, para apartar tu c\u00f3lera.<\/p>\n<p><span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/p>\n<p><\/span><\/p>\n<p><span class=\"Tit_Lectura\">Sal 30, 5 6. 14. 15 16<\/span><br \/><span class=\"Tit_Negro_Normal\">R. S\u00e1lvame, Se\u00f1or, por tu misericordia.<\/span><\/p>\n<p>S\u00e1came de la red que me han tendido, <br \/>porque t\u00fa eres mi amparo. <br \/>A tus manos encomiendo mi esp\u00edritu: <br \/>t\u00fa, el Dios leal, me librar\u00e1s, R.<\/p>\n<p>Oigo el cuchicheo de la gente, <br \/>y todo me da miedo; <br \/>se conjuran contra m\u00ed <br \/>y traman quitarme la vida. R.<\/p>\n<p>Pero yo conf\u00edo en ti, Se\u00f1or, <br \/>te digo: \u00abT\u00fa eres mi Dios.\u00bb <br \/>En tu mano est\u00e1n mis azares: <br \/>l\u00edbrame de mis enemigos que me persiguen. R. <\/p>\n<p><span class=\"Tit_Lectura\">Vers\u00edculo  Jn 8, 12b<\/span><br \/><span class=\"Tit_Negro_Normal\"><\/span><\/p>\n<p>V: Yo soy la luz del mundo - dice el Se\u00f1or -;<br \/>el que me sigue tendr\u00e1 la luz de la vida. <\/p>\n<p><span class=\"Tit_Negro_Cur\">EVANGELIO<\/span><br \/> \t\t\t\t\t\t<span class=\"Tit_Lectura\">Lo condenar\u00e1n a muerte.<\/span><br \/><span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio seg\u00fan san Mateo 20, 17-28<\/span><\/p>\n<p>En aquel tiempo, subiendo Jes\u00fas a Jerusal\u00e9n, tomando aparte a los Doce, les dijo por el camino:<\/p>\n<p>\u00abMirad, estamos subiendo a Jerusal\u00e9n, y el Hijo del hombre va a ser entregado a los sumos sacerdotes y a los escribas, y lo condenar\u00e1n a muerte y lo entregar\u00e1n a los gentiles, para que se burlen de \u00e9l, lo azoten y lo crucifiquen; y al tercer d\u00eda resucitar\u00e1\u00bb.<\/p>\n<p>Entonces se le acerc\u00f3 la madre de los hijos de Zebedeo con sus hijos y se postr\u00f3 para hacerle una petici\u00f3n. <\/p>\n<p>\u00c9l le pregunt\u00f3:<\/p>\n<p>\u00ab\u00bfQu\u00e9 deseas?\u00bb.<\/p>\n<p>Ella contest\u00f3:<\/p>\n<p>\u00abOrdena que estos dos hijos m\u00edos se sienten en tu reino, uno a tu derecha y el otro a tu izquierda\u00bb<\/p>\n<p>Pero Jes\u00fas replic\u00f3:<\/p>\n<p>\u00abNo sab\u00e9is lo que ped\u00eds. \u00bfPod\u00e9is beber el c\u00e1liz que yo he de beber?\u00bb<\/p>\n<p>Contestaron:<\/p>\n<p>\u00abLo somos.\u00bb<\/p>\n<p>\u00c9l les dijo:<\/p>\n<p>\u00abMi c\u00e1liz lo beber\u00e9is; pero sentarse a mi derecha o a mi izquierda no me toca a m\u00ed concederlo, es para aquellos para quienes lo tiene reservado mi Padre\u00bb.<\/p>\n<p>Los otros diez, al o\u00edr aquello, se indignaron contra los dos hermanos. Y llam\u00e1ndolos, Jes\u00fas les dijo:<\/p>\n<p>\u00abSab\u00e9is que los jefes de los pueblos los tiranizan y que los grandes los oprimen. No ser\u00e1 as\u00ed entre vosotros: el que quiera ser grande entre vosotros, que sea vuestro servidor, y el que quiera ser primero entre vosotros, que sea vuestro esclavo.<\/p>\n<p>Igual que el Hijo del hombre no ha venido a ser servido sino a servir y a dar su vida en rescate por muchos\u00bb.<\/p>\n<p class=\"Tit_Negro_Normal\">Palabra del Se\u00f1or.<\/p>\n","post_title":"16\/03\/2022 - Mi\u00e9rcoles de la 2\u00aa semana de Cuaresma.","post_excerpt":"","post_status":"publish","comment_status":"open","ping_status":"open","post_password":"","post_name":"16-03-2022-miercoles-de-la-2a-semana-de-cuaresma","to_ping":"","pinged":"","post_modified":"2022-03-16 00:00:00","post_modified_gmt":"2022-03-15 23:00:00","post_content_filtered":"","post_parent":"0","guid":"https:\/\/oracionyliturgia.archimadrid.org\/?p=49445","menu_order":"0","post_type":"post","post_mime_type":"","comment_count":"3"}]`,
			code:        http.StatusOK,
			expected: &Gospel{
				Day:       "16/03/2022 - Miércoles de la 2ª semana de Cuaresma.",
				Title:     "Sal 30, 5 6. 14. 15 16",
				Reference: "R. Sálvame, Señor, por tu misericordia.",
				Content:   "Sácame de la red que me han tendido, porque tú eres mi amparo. A tus manos encomiendo mi espíritu: tú, el Dios leal, me librarás, R.\nOigo el cuchicheo de la gente, y todo me da miedo; se conjuran contra mí y traman quitarme la vida. R.\nPero yo confío en ti, Señor, te digo: «Tú eres mi Dios.» En tu mano están mis azares: líbrame de mis enemigos que me persiguen. R. ",
			},
			errorExpected: false,
		},
		{
			name:        "Error in archimadrid server",
			day:         time.Date(2022, time.March, 16, 0, 0, 0, 0, time.UTC),
			regexString: `(Palabra\sde\sDios\..*<p>)<span.*?\sR.\s`,
			cachePrefix: "psalm ",
			psalm:       false,
			response:    "",
			code:        http.StatusInternalServerError,
			expected: &Gospel{
				Day: "2022-03-16",
			},
			errorExpected: true,
		},
		{
			name:          "Empty First Lecture",
			day:           time.Date(2022, time.March, 10, 0, 0, 0, 0, time.UTC),
			regexString:   `(PRIMERA\sLECTURA).*?Palabra de Dios\.`,
			cachePrefix:   "first lecture ",
			psalm:         false,
			response:      `[{"ID":"49439","post_author":"6","post_date":"2022-03-10 00:00:00","post_date_gmt":"2022-03-09 23:00:00","post_content":"<span class=\"Tit_Negro_Cur\">PRIMERA LECTURA<\/span>\r\n<span class=\"Tit_Lectura\">No tengo m\u00e1s defensor que t\u00fa.<\/span>\r\n<span class=\"Tit_Negro_Normal\">Lectura del libro de Ester 4, 17k. l-z<\/span>\r\n\r\nEn aquellos d\u00edas, la reina Ester, presa de un temor mortal, se refugi\u00f3 en el Se\u00f1or.\r\n\r\nY se postro en tierra con sus doncellas desde la ma\u00f1ana a la tarde, diciendo:\r\n\r\n\u00ab\u00a1Bendita seas, Dios de Abrah\u00e1n, Dios de Isaac y Dios de Jacob! Ven en mi ayuda, que estoy sola y no tengo otro socorro fuera de ti, Se\u00f1or, por que me acecha un gran peligro.\r\n\r\nYo he escuchado en los libros de mis antepasados, Se\u00f1or, que t\u00fa libras siempre a los que cumplen tu voluntad. Ahora, Se\u00f1or, Dios m\u00edo, ay\u00fadame, que estoy sola y no tengo a nadie fuera de ti. Ahora, ven en mi ayuda, pues estoy hu\u00e9rfana, y pon en mis labios una palabra oportuna de lante del le\u00f3n, y hazme grata a sus ojos. Cambia su coraz\u00f3n para que aborrezca al que nos ataca, para su ruina y la de cuantos est\u00e1n de acuerdo con \u00e9l.\r\n\r\nL\u00edbranos de la mano de nuestros enemigos, cambia nuestro luto en gozo y nuestros sufrimientos en salvaci\u00f3n\u00bb.\r\n\r\n<span class=\"Tit_Negro_Normal\">Palabra de Dios.<\/span>\r\n\r\n&nbsp;\r\n\r\n<span class=\"Tit_Lectura\">Sal 137, 1-2a. 2bc y 3. 7c-8<\/span>\r\n<span class=\"Tit_Negro_Normal\">R. Cuando te invoqu\u00e9, me escuchaste, Se\u00f1or.<\/span>\r\n\r\nTe doy gracias, Se\u00f1or, de todo coraz\u00f3n;\r\nporque escuchaste las palabras de mi boca;\r\ndelante de los \u00e1ngeles ta\u00f1er\u00e9 para ti;\r\nme postrar\u00e9 hacia tu santuario. R.\r\n\r\nDar\u00e9 gracias a tu nombre,\r\npor tu misericordia y tu lealtad;\r\nporque tu promesa supera tu fama.\r\nCuando te invoqu\u00e9, me escuchaste,\r\nacreciste el valor en mi alma. R.\r\n\r\nTu derecha me salva.\r\nEl Se\u00f1or completar\u00e1 sus favores conmigo:\r\nSe\u00f1or, tu misericordia es eterna,\r\nno abandones la obra de tus manos. R.\r\n\r\n<span class=\"Tit_Lectura\">Vers\u00edculo Sal 50, 12a. 14a<\/span>\r\n\r\nV: Oh, Dios, crea en m\u00ed un coraz\u00f3n puro;\r\ny devu\u00e9lveme la alegr\u00eda de tu salvaci\u00f3n.\r\n\r\n<span class=\"Tit_Negro_Cur\">EVANGELIO<\/span>\r\n<span class=\"Tit_Lectura\">Todo el que pide recibe.<\/span>\r\n<span class=\"Tit_Negro_Normal\">Lectura del santo Evangelio seg\u00fan san Mateo 7, 7-12<\/span>\r\n\r\nEn aquel tiempo, dijo Jes\u00fas a sus disc\u00edpulos:\r\n\r\n\u00abPedid y se os dar\u00e1, buscad y encontrar\u00e9is, llamad y se os abrir\u00e1; porque todo el que pide recibe, quien busca encuentra y al que llama se le abre.\r\n\r\nSi a alguno de vosotros le pide su hijo pan, \u00bfle va a dar una piedra?; y si le pide pescado, \u00bfle dar\u00e1 una serpiente? Pues si vosotros, aun siendo malos, sab\u00e9is dar cosas buenas a vuestros hijos, \u00a1cu\u00e1nto m\u00e1s vuestro Padre que est\u00e1 en los cielos dar\u00e1 cosas buenas a los que le piden!\r\n\r\nAs\u00ed, pues, todo lo que dese\u00e1is que los dem\u00e1s hagan con vosotros, hacedlo vosotros con ellos; pues esta es la Ley y los profetas\u00bb.\r\n<p class=\"Tit_Negro_Normal\">Palabra del Se\u00f1or.<\/p>","post_title":"10\/03\/2022 - Jueves de la 1\u00aa semana de Cuaresma.","post_excerpt":"","post_status":"publish","comment_status":"open","ping_status":"open","post_password":"","post_name":"10-03-2022-jueves-de-la-1a-semana-de-cuaresma","to_ping":"","pinged":"","post_modified":"2022-02-25 13:27:20","post_modified_gmt":"2022-02-25 12:27:20","post_content_filtered":"","post_parent":"0","guid":"https:\/\/oracionyliturgia.archimadrid.org\/?p=49439","menu_order":"0","post_type":"post","post_mime_type":"","comment_count":"2"}]`,
			code:          http.StatusOK,
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.code)
				w.Write([]byte(test.response))
			}))
			defer server.Close()

			client := NewClient(SetURL(server.URL))

			if test.cache != nil {
				client.Set(test.key, test.cache)
			}
			actual, err := client.getGospelOrLecture(
				context.TODO(),
				test.day,
				test.regexString,
				test.cachePrefix,
				test.psalm,
			)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.EqualValues(tt, test.expected, actual)

			today := test.day.Format("2006-01-02")
			if test.response != "" {
				object, err := client.Get(ResponsePrefix + today)
				_, ok := object.(*gospelResponse)
				assert.NoError(tt, err)
				assert.True(tt, ok)
			}

			object, err := client.Get(test.cachePrefix + today)
			actual, ok := object.(*Gospel)
			assert.NoError(tt, err)
			assert.True(tt, ok)
			assert.EqualValues(tt, test.expected, actual)
		})
	}
}
