package guess

import (
	"ga/forces/models"
	"math"
	"path/filepath"
	"testing"

	"github.com/go-test/deep"
)

const PathToH2CODir = "../testfiles/h2co/4th"

var fort15 = models.Chromosome{
	0.0177072648, 0.0000000000, 0.0000000000,
	-0.0546136978, 0.0000000000, 0.0000000000,
	0.0184532165, 0.0000000000, 0.0000000000,
	0.0184532165, 0.0000000000, 0.0000000000,
	0.0000000000, 0.0725713528, 0.0000000000,
	0.0000000000, -0.1045740584, 0.0000000000,
	0.0000000000, 0.0160013528, -0.0367959445,
	0.0000000000, 0.0160013528, 0.0367959445,
	0.0000000000, 0.0000000000, 0.8391376763,
	0.0000000000, 0.0000000000, -0.7564487448,
	0.0000000000, -0.0151845601, -0.0413444658,
	0.0000000000, 0.0151845601, -0.0413444658,
	-0.0546136978, 0.0000000000, 0.0000000000,
	0.1684425025, 0.0000000000, 0.0000000000,
	-0.0569144023, 0.0000000000, 0.0000000000,
	-0.0569144023, 0.0000000000, 0.0000000000,
	0.0000000000, -0.1045740584, 0.0000000000,
	0.0000000000, 0.5951947650, 0.0000000000,
	0.0000000000, -0.2453103533, -0.0841298010,
	0.0000000000, -0.2453103533, 0.0841298010,
	0.0000000000, 0.0000000000, -0.7564487448,
	0.0000000000, 0.0000000000, 0.9541392898,
	0.0000000000, -0.0820167905, -0.0988452725,
	0.0000000000, 0.0820167905, -0.0988452725,
	0.0184532165, 0.0000000000, 0.0000000000,
	-0.0569144023, 0.0000000000, 0.0000000000,
	0.0192305929, 0.0000000000, 0.0000000000,
	0.0192305929, 0.0000000000, 0.0000000000,
	0.0000000000, 0.0160013528, -0.0151845601,
	0.0000000000, -0.2453103533, -0.0820167905,
	0.0000000000, 0.2472859793, 0.1090635481,
	0.0000000000, -0.0179769788, -0.0118621974,
	0.0000000000, -0.0367959445, -0.0413444658,
	0.0000000000, -0.0841298010, -0.0988452725,
	0.0000000000, 0.1090635481, 0.1310698454,
	0.0000000000, 0.0118621974, 0.0091198929,
	0.0184532165, 0.0000000000, 0.0000000000,
	-0.0569144023, 0.0000000000, 0.0000000000,
	0.0192305929, 0.0000000000, 0.0000000000,
	0.0192305929, 0.0000000000, 0.0000000000,
	0.0000000000, 0.0160013528, 0.0151845601,
	0.0000000000, -0.2453103533, 0.0820167905,
	0.0000000000, -0.0179769788, 0.0118621974,
	0.0000000000, 0.2472859793, -0.1090635481,
	0.0000000000, 0.0367959445, -0.0413444658,
	0.0000000000, 0.0841298010, -0.0988452725,
	0.0000000000, -0.0118621974, 0.0091198929,
	0.0000000000, -0.1090635481, 0.1310698454,
}

var fort30 = models.Chromosome{
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, -0.3339623021, 0.000000000,
	-0.2438654704, 0.0000000000, 0.000000000,
	2.9961819173, 0.0000000000, 0.000000000,
	0.0000000000, 0.2789775987, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.2187387642, 0.0000000000, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	0.2237211528, 0.0000000000, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, 0.0000000000, -0.296433575,
	0.0000000000, 0.0000000000, 0.294885113,
	0.0000000000, 0.2466385228, 0.000000000,
	0.0000000000, -2.8650309731, -0.267801574,
	0.0000000000, 0.0000000000, 0.107592947,
	0.0000000000, -0.2311064601, 0.000000000,
	0.0000000000, 1.2178420205, 0.000000000,
	0.0000000000, 2.7976692350, 0.000000000,
	0.0000000000, -2.6664483296, 0.000000000,
	0.0154778256, 0.0000000000, 0.027492351,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0140944245, -0.0301194172, 0.000000000,
	-0.0140944245, 0.0000000000, 0.000000000,
	0.1203919961, 0.0000000000, -0.013541769,
	0.0000000000, 0.0000000000, 0.080104313,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0013834011, -0.0029702161, 0.000000000,
	-0.1062975717, -0.0564192049, 0.000000000,
	0.0037874874, 0.0000000000, 0.009287872,
	0.0000000000, 0.0100721588, -0.047080878,
	0.0022019126, 0.0000000000, 0.000000000,
	0.0681979985, 0.0000000000, 0.009304630,
	0.0363562115, 0.0000000000, -0.562672319,
	0.0000000000, -0.0077660313, 0.050276709,
	0.0000000000, -0.4933677802, -0.272015063,
	-0.0036864005, 0.0000000000, 0.000000000,
	-0.0883487413, 0.0000000000, 0.000000000,
	0.1103547651, 0.0000000000, -0.018592503,
	-0.0333944399, 0.0000000000, 0.553367689,
	0.4552736331, 0.0000000000, -0.556264190,
	0.0195385946, 0.0000000000, -0.001386526,
	0.0000000000, -0.0532131260, -0.065575472,
	-0.0055880120, 0.0000000000, 0.000000000,
	0.0555729082, 0.0000000000, 0.003692653,
	0.0624629058, 0.0000000000, -0.460704222,
	0.0000000000, 0.0204969933, 0.033680869,
	0.0000000000, -0.2570960984, -0.065610452,
	-0.0173558741, 0.0000000000, 0.000000000,
	-0.0487264867, 0.0000000000, 0.000000000,
	0.0638093319, 0.0000000000, 0.001801553,
	-0.0062228057, 0.0000000000, 0.440135566,
	0.2291687294, 0.0000000000, -0.425523995,
	0.0000000000, 0.0327161327, 0.032650122,
	0.0000000000, 0.1946331926, 0.028913924,
	0.0000000000, -0.2162052378, -0.060433976,
	0.0000000000, -0.0154778256, 0.000000000,
	0.0274923517, 0.0000000000, 0.000000000,
	0.0000000000, 0.0140944245, -0.030119417,
	0.0000000000, 0.0140944245, 0.000000000,
	0.0000000000, -0.1203919961, 0.000000000,
	-0.0135417691, 0.0000000000, 0.000000000,
	0.0801043134, 0.0000000000, 0.000000000,
	0.0000000000, 0.0000000000, 0.005597281,
	0.0000000000, 0.0000000000, -0.010143339,
	0.0000000000, -0.0023029994, 0.000000000,
	0.0000000000, 0.0179488303, 0.000000000,
	0.0000000000, -0.0183196233, 0.000000000,
	0.0034052916, 0.0000000000, 0.000000000,
	-0.0012584095, 0.0000000000, 0.000000000,
	0.0022730289, 0.0000000000, 0.000000000,
	0.0000000000, 0.0013834011, -0.002970216,
	0.0000000000, 0.1062975717, -0.056419204,
	0.0000000000, 0.0026737924, -0.004419911,
	0.0000000000, -0.0037874874, 0.000000000,
	-0.0092878729, 0.0000000000, 0.010072158,
	0.0470808782, -0.0022019126, 0.000000000,
	0.0000000000, -0.0681979985, 0.000000000,
	-0.0093046301, 0.0363562115, 0.000000000,
	0.5626723192, 0.0000000000, -0.007766031,
	-0.0502767099, 0.0000000000, -0.493367780,
	0.2720150636, 0.0023029994, 0.000000000,
	0.0000000000, -0.0179488303, 0.000000000,
	0.0000000000, -0.0026737924, 0.000000000,
	0.0000000000, -0.0130339303, 0.000000000,
	0.0000000000, 0.0458601784, 0.000000000,
	0.0214890042, 0.0000000000, -0.004107681,
	-0.0030269740, 0.0000000000, 0.016876002,
	0.0074303757, 0.0000000000, -0.016413124,
	-0.0111440876, 0.0036864005, 0.000000000,
	0.0000000000, 0.0883487413, 0.000000000,
	0.0000000000, 0.0183196233, 0.000000000,
	0.0000000000, -0.1103547651, 0.000000000,
	0.0185925031, -0.0333944399, 0.000000000,
	-0.5533676891, 0.4552736331, 0.000000000,
	-0.0214890042, 0.0036448026, 0.000000000,
	0.5562641902, 0.0195385946, 0.000000000,
	-0.0013865262, 0.0000000000, 0.053213126,
	-0.0655754721, -0.0055880120, 0.000000000,
	0.0000000000, 0.0555729082, 0.000000000,
	0.0036926536, -0.0624629058, 0.000000000,
	-0.4607042224, 0.0000000000, -0.020496993,
	0.0336808690, 0.0000000000, 0.257096098,
	-0.0656104527, 0.0034052916, 0.000000000,
	0.0000000000, -0.0012584095, 0.000000000,
	0.0000000000, -0.0044199110, 0.000000000,
	-0.0041076811, 0.0030269740, 0.000000000,
	0.0168760025, -0.0074303757, 0.000000000,
	0.0036448026, 0.0000000000, 0.000000000,
	-0.0007555192, 0.0000000000, -0.000000000,
	0.0030156592, 0.0000000000, -0.006740685,
	-0.0011300700, -0.0173558741, 0.000000000,
	0.0000000000, -0.0487264867, 0.000000000,
	0.0000000000, 0.0022730289, 0.000000000,
	0.0000000000, 0.0638093319, 0.000000000,
	0.0018015537, 0.0062228057, 0.000000000,
	0.4401355662, -0.2291687294, 0.000000000,
	-0.0164131240, 0.0067406859, 0.000000000,
	-0.4255239958, 0.0000000000, -0.032716132,
	0.0326501223, 0.0000000000, -0.194633192,
	0.0289139245, 0.0000000000, 0.011144087,
	-0.0011300700, 0.0000000000, 0.216205237,
	-0.0604339768,
}

var fort40 = models.Chromosome{
	0.4084362901, 0.0000000000, 0.102262228,
	0.0000000000, 0.2410739049, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	-1.5746834458, 0.0000000000, -1.413580580,
	0.0000000000, 0.0000000000, 8.973451013,
	-0.3436800333, 0.0000000000, -0.092432585,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, 1.5065724802, 0.000000000,
	0.0000000000, 0.2802309778, 0.000000000,
	0.0595800249, 0.0000000000, 0.000000000,
	-1.5169611495, -0.1745966400, 0.000000000,
	0.0000000000, 0.2433785138, 0.000000000,
	-0.0900091366, 0.0000000000, -0.230914333,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, 1.4307549884, 0.000000000,
	0.0000000000, 0.0730592106, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, -0.0140288337, -0.000000000,
	0.0000000000, 0.1134270146, 0.000000000,
	0.2302803198, 0.0000000000, 0.000000000,
	-1.5995569071, -0.0788037642, 0.000000000,
	0.0000000000, -0.8704415157, 0.000000000,
	-0.2293319392, 0.0000000000, 0.000000000,
	1.6980345796, 0.0000000000, 0.000000000,
	0.0000000000, 0.0000000000, 1.497027564,
	0.0000000000, 1.4221054786, 0.000000000,
	0.0000000000, -8.8146091926, 0.000000000,
	-0.0000000000, 0.0000000000, -1.472977713,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, 1.5402320140, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0000000000, -1.4480393224, 0.000000000,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0000000000, 0.0000000000, 0.000000000,
	1.5576686928, 0.0000000000, 0.000000000,
	-1.4602488920, 0.0000000000, -1.448684964,
	0.0000000000, 0.0000000000, 8.722582547,
	1.4561141894, 0.0000000000, 0.000000000,
	-1.8563016669, 0.0000000000, 1.527427803,
	0.0000000000, 0.0000000000, -0.284750562,
	0.0000000000, 0.0000000000, -8.629589563,
	0.0000000000, 0.0000000000, 7.730299131,
	-0.0323781284, 0.0000000000, -0.004914821,
	0.0000000000, 0.0000000000, 0.025417410,
	0.0000000000, 0.0340554828, 0.000000000,
	0.0000000000, 0.0317245278, 0.000000000,
	0.0164262801, 0.0000000000, -0.015588514,
	0.0051943346, -0.0528171689, 0.000000000,
	0.0000000000, -0.0343909369, 0.000000000,
	0.0084749630, 0.0000000000, -0.015588514,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0295151884, -0.0219186499, 0.000000000,
	-0.0173116252, 0.0000000000, 0.000000000,
	0.4746226399, 0.0000000000, 0.000000000,
	-0.0136109491, 0.0000000000, -0.012024925,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0103521375, -0.0336271501, 0.000000000,
	0.0103521375, 0.0000000000, 0.000000000,
	0.3132147473, 0.0000000000, 0.002067351,
	0.0000000000, 0.0000000000, 0.200093738,
	0.0000000000, 0.0000000000, 0.010128261,
	0.0000000000, -0.0118684726, 0.000000000,
	-0.0098288960, -0.0329069962, 0.003215090,
	0.0000000000, 0.0000000000, 0.130771625,
	0.0000000000, 0.0210974785, 0.037507164,
	0.0000000000, -0.4553202908, 0.000000000,
	0.0032588115, 0.0370097081, 0.000000000,
	-0.3235668848, -0.1778192981, -0.009141906,
	0.0000000000, 0.0000000000, -0.146829699,
	0.0000000000, 0.0000000000, 0.171383061,
	0.0000000000, -0.0061265459, 0.000000000,
	-0.0050797860, 0.0157691088, 0.000000000,
	0.0196611387, 0.0000000000, -0.008587203,
	-0.0979823845, 0.0000000000, 0.009686687,
	0.0000000000, -0.0021261777, 0.000000000,
	0.0000000000, 0.0000000000, -0.022775595,
	-0.0373978802, 0.0000000000, -0.011708939,
	0.0000000000, 0.0003170066, 0.000000000,
	0.0099311564, 0.0844009593, 0.002872276,
	0.0000000000, 0.0000000000, 0.442235174,
	0.0000000000, -0.0004741903, -0.015337692,
	0.0000000000, -0.7343513202, -0.018106593,
	0.0000000000, -0.0254507553, 0.000000000,
	0.0129669219, 0.0730187319, 0.011889777,
	0.0000000000, 0.0000000000, 0.274391122,
	0.0000000000, 0.0185872217, -0.054814685,
	0.0000000000, -1.2540126307, 0.000000000,
	-0.0393714197, -0.0561284685, 0.000000000,
	-0.6213386203, 0.0192230760, 0.000000000,
	-0.0045948277, 0.0000000000, -0.011735913,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0125244254, 0.0385156111, 0.000000000,
	0.0103983023, 0.0000000000, 0.000000000,
	-0.4344939485, 0.0000000000, 0.004737814,
	0.0000000000, 0.0000000000, -0.304923892,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0077797172, -0.0219728918, 0.000000000,
	0.4230716577, 0.2944478031, 0.000000000,
	0.0123616475, 0.0000000000, -0.002840773,
	0.0000000000, -0.0295922951, -0.052142709,
	-0.0075933903, 0.0000000000, 0.000000000,
	-0.3928439783, 0.0000000000, 0.001968731,
	0.0054065360, 0.0000000000, 0.753637608,
	0.0000000000, 0.0068635336, 0.012091986,
	0.0000000000, 1.2354254091, 0.668667852,
	-0.0055493428, 0.0000000000, 0.000000000,
	0.4061789359, 0.0000000000, 0.000000000,
	-0.4195085194, 0.0000000000, 0.003768043,
	0.0161071288, 0.0000000000, -0.765918163,
	-1.2044277948, 0.0000000000, 0.783371515,
	0.0000000000, 0.0246696555, 0.000000000,
	0.0482237278, 0.0388279407, 0.000000000,
	-0.0042624488, 0.0000000000, -0.079084657,
	-0.0794209105, 0.0000000000, -0.013282248,
	0.0000000000, -0.0167973831, 0.000000000,
	0.0000000000, 0.0000000000, 0.009338341,
	-0.0116354323, 0.0000000000, -0.032531702,
	0.0000000000, -0.0471128217, 0.000000000,
	0.0086421670, 0.0664593336, 0.030338440,
	0.0000000000, 0.0000000000, 0.251571917,
	0.0000000000, 0.0656383142, 0.020944107,
	0.0000000000, -1.3363896797, -0.018389336,
	0.0000000000, 0.0132897427, 0.000000000,
	0.0344051304, 0.0460133224, 0.008431762,
	0.0000000000, 0.0000000000, 0.158034826,
	0.0000000000, -0.0396942406, -0.040679130,
	0.0000000000, -0.6364590649, 0.000000000,
	-0.0123043091, -0.0464964922, 0.000000000,
	0.0736621006, 0.4496452159, 0.000000000,
	-0.0115969342, 0.0000000000, -0.026027347,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0045901420, 0.0318379520, 0.000000000,
	0.0037148194, 0.0000000000, 0.000000000,
	-0.2866032276, 0.0000000000, 0.010761237,
	0.0000000000, 0.0000000000, -0.155763226,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0089079392, -0.0046608312, 0.000000000,
	0.2793581265, 0.1371522760, 0.000000000,
	0.0050997656, 0.0000000000, 0.002339355,
	0.0000000000, 0.0069498093, 0.018794488,
	-0.0134098957, 0.0000000000, 0.000000000,
	-0.2489517504, 0.0000000000, -0.023521935,
	-0.0486919075, 0.0000000000, 1.270050844,
	0.0000000000, 0.0143461100, -0.005308131,
	0.0000000000, 0.6978457422, -0.012226199,
	0.0095596081, 0.0000000000, 0.000000000,
	0.2690947454, 0.0000000000, 0.000000000,
	-0.2773740514, 0.0000000000, 0.021864836,
	0.0447844636, 0.0000000000, -1.231597770,
	-0.7012928704, 0.0000000000, 1.187034977,
	-0.0225899393, 0.0000000000, -0.003224564,
	0.0000000000, 0.0446795265, 0.028152172,
	0.0080830758, 0.0000000000, 0.000000000,
	-0.1404695206, 0.0000000000, 0.019329038,
	-0.0257802037, 0.0000000000, 0.643959179,
	0.0000000000, -0.0221008213, 0.017639276,
	0.0000000000, -0.0329829706, -0.434392066,
	0.0164830501, 0.0000000000, 0.000000000,
	0.1211016099, 0.0000000000, 0.000000000,
	-0.1297864534, 0.0000000000, -0.018335766,
	-0.0164567853, 0.0000000000, -0.657514428,
	0.0242432979, 0.0000000000, 0.652515430,
	0.0000000000, -0.0223791973, -0.051741794,
	0.0000000000, 0.0524717981, 0.423796408,
	0.0000000000, -0.0066098471, -0.363902766,
	-0.0323781284, 0.0000000000, -0.004914821,
	0.0000000000, 0.0000000000, -0.025417410,
	0.0000000000, 0.0340554828, 0.000000000,
	0.0000000000, 0.0317245278, 0.000000000,
	0.0164262801, 0.0000000000, 0.015588514,
	0.0051943346, -0.0528171689, 0.000000000,
	0.0000000000, -0.0343909369, 0.000000000,
	0.0084749630, 0.0000000000, 0.015588514,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0295151884, 0.0219186499, 0.000000000,
	-0.0173116252, 0.0000000000, 0.000000000,
	0.4746226399, 0.0000000000, 0.000000000,
	0.0136109491, 0.0000000000, -0.012024925,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0103521375, -0.0336271501, 0.000000000,
	-0.0103521375, 0.0000000000, 0.000000000,
	-0.3132147473, 0.0000000000, 0.002067351,
	0.0000000000, 0.0000000000, 0.200093738,
	0.0000000000, 0.0000000000, -0.009474661,
	0.0000000000, 0.0003570142, 0.000000000,
	0.0000000000, -0.0063428212, 0.017877550,
	0.0000000000, 0.0000000000, -0.043563519,
	0.0000000000, -0.0000572531, -0.000000000,
	0.0000000000, -0.0019907239, 0.000000000,
	0.0000000000, 0.0086423672, 0.000000000,
	0.0000000000, -0.0243417919, -0.004201444,
	0.0000000000, 0.0000000000, 0.012842984,
	0.0000000000, 0.0000000000, -0.015411455,
	0.0000000000, 0.0010346863, 0.000000000,
	-0.0019070176, 0.0000000000, 0.000000000,
	0.0000000000, 0.0005644829, 0.001008446,
	0.0000000000, -0.0015616401, 0.000000000,
	0.0000000000, -0.0106135030, 0.000000000,
	0.0014790024, 0.0000000000, 0.000000000,
	0.0186429928, 0.0000000000, 0.000000000,
	0.0000000000, -0.0001498805, -0.004806805,
	0.0000000000, 0.0010239885, 0.005738274,
	0.0000000000, 0.0007810856, 0.000000000,
	0.0000000000, -0.0057415672, 0.000000000,
	0.0000000000, 0.0188789263, 0.000000000,
	0.0000000000, 0.0002095272, 0.000000000,
	0.0039967903, 0.0000000000, 0.000000000,
	0.0000000000, -0.0006462350, -0.003405136,
	0.0000000000, -0.0015215576, 0.000000000,
	0.0000000000, 0.0046928698, 0.000000000,
	-0.0008036637, 0.0000000000, 0.000000000,
	-0.0107033626, 0.0000000000, 0.000000000,
	0.0000000000, -0.0019011469, -0.001149773,
	0.0000000000, 0.0035302818, 0.007849712,
	0.0000000000, -0.0012494780, 0.000000000,
	0.0000000000, -0.0067330993, 0.000000000,
	0.0000000000, -0.0012803020, 0.000000000,
	-0.0019761866, 0.0000000000, 0.000000000,
	0.0112848350, 0.0000000000, 0.000000000,
	-0.0077982066, 0.0000000000, 0.000000000,
	0.0101282616, 0.0000000000, -0.011868472,
	0.0000000000, 0.0098288960, -0.032906996,
	0.0032150902, 0.0000000000, 0.000000000,
	0.1307716250, 0.0000000000, 0.021097478,
	-0.0375071644, 0.0000000000, -0.455320290,
	0.0000000000, -0.0032588115, 0.037009708,
	0.0000000000, 0.3235668848, -0.177819298,
	-0.0042014449, 0.0000000000, 0.000000000,
	0.0128429842, 0.0000000000, 0.000000000,
	0.0067699161, 0.0000000000, -0.001449288,
	0.0057053766, 0.0000000000, 0.011151154,
	-0.0258602702, 0.0000000000, -0.013918444,
	0.0000000000, 0.0023378547, 0.000558119,
	0.0000000000, -0.0067015939, 0.003657313,
	0.0000000000, 0.0092628794, -0.001510441,
	-0.0091419068, 0.0000000000, 0.000000000,
	-0.1468296994, 0.0000000000, 0.000000000,
	-0.0154114554, 0.0000000000, 0.000000000,
	0.1713830616, 0.0000000000, -0.006126545,
	0.0000000000, -0.0050797860, -0.015769108,
	0.0000000000, -0.0196611387, 0.000000000,
	-0.0085872037, 0.0979823845, 0.000000000,
	0.0096866872, 0.0000000000, 0.002126177,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0227755956, 0.0373978802, 0.000000000,
	-0.0117089390, 0.0000000000, 0.000317006,
	0.0000000000, -0.0099311564, 0.084400959,
	0.0028722768, 0.0000000000, 0.000000000,
	0.4422351747, 0.0000000000, -0.000474190,
	0.0153376924, 0.0000000000, -0.734351320,
	0.0181065934, 0.0000000000, 0.025450755,
	0.0000000000, 0.0129669219, -0.073018731,
	-0.0118897770, 0.0000000000, 0.000000000,
	-0.2743911222, 0.0000000000, -0.018587221,
	-0.0548146852, 0.0000000000, 1.254012630,
	0.0000000000, -0.0393714197, 0.056128468,
	0.0000000000, -0.6213386203, -0.019223076,
	0.0000000000, 0.0010346863, 0.000000000,
	0.0019070176, 0.0000000000, 0.000000000,
	0.0000000000, 0.0005644829, -0.001008446,
	0.0000000000, -0.0015616401, 0.000000000,
	0.0000000000, -0.0106135030, 0.000000000,
	-0.0014790024, 0.0000000000, 0.000000000,
	-0.0186429928, 0.0000000000, 0.000000000,
	0.0000000000, -0.0014492887, -0.005705376,
	0.0000000000, 0.0111511546, 0.025860270,
	0.0000000000, 0.0054738374, 0.000000000,
	0.0076035524, 0.0000000000, -0.000000000,
	-0.0236710457, -0.0049655738, 0.000000000,
	0.0000000000, -0.0266156008, 0.000000000,
	-0.0018115481, -0.0000000000, 0.000000000,
	-0.0188120984, 0.0000000000, 0.000000000,
	0.0297557765, 0.0000000000, 0.000000000,
	-0.0079578129, -0.0002541318, 0.000000000,
	0.0000000000, 0.0157905873, 0.000000000,
	0.0000000000, 0.0042165789, 0.000000000,
	-0.0028960022, 0.0080786303, 0.000000000,
	0.0103118232, -0.0378611479, 0.000000000,
	-0.0212213951, 0.0027622811, 0.000000000,
	-0.0034502614, 0.0000000000, -0.011329527,
	-0.0061691646, -0.0036462963, 0.000000000,
	0.0000000000, -0.0119585084, 0.000000000,
	0.0049964428, 0.0191056334, 0.000000000,
	0.0007005211, 0.0000000000, 0.012058387,
	0.0115821314, 0.0000000000, -0.021692436,
	-0.0491315920, -0.0016774932, 0.000000000,
	0.0000000000, 0.0129183402, 0.000000000,
	0.0000000000, -0.0108920142, 0.000000000,
	-0.0006822567, -0.0030423654, 0.000000000,
	-0.0149311382, -0.0108989818, 0.000000000,
	0.0226979563, 0.0000000000, 0.002231292,
	-0.0024425376, 0.0000000000, -0.005773790,
	0.0308404941, 0.0000000000, 0.023334763,
	-0.0234827537, 0.0000000000, -0.004594827,
	0.0000000000, 0.0117359136, 0.000000000,
	0.0000000000, 0.0000000000, 0.012524425,
	-0.0385156111, 0.0000000000, 0.010398302,
	0.0000000000, 0.0000000000, -0.434493948,
	0.0000000000, -0.0047378140, 0.000000000,
	0.0000000000, 0.3049238920, 0.000000000,
	0.0000000000, 0.0000000000, -0.000149880,
	0.0048068057, 0.0000000000, 0.001023988,
	-0.0057382749, 0.0000000000, -0.000254131,
	0.0000000000, 0.0000000000, 0.015790587,
	0.0000000000, 0.0000000000, -0.019753034,
	0.0000000000, 0.0025615085, 0.000000000,
	0.0000000000, 0.0026864645, 0.000000000,
	0.0000000000, -0.0003488328, 0.000000000,
	0.0000000000, 0.0000000000, -0.007779717,
	0.0219728918, 0.0000000000, 0.423071657,
	-0.2944478031, 0.0000000000, 0.004216578,
	-0.0048991401, 0.0000000000, 0.012361647,
	0.0000000000, -0.0028407730, 0.000000000,
	0.0295922951, -0.0521427099, -0.007593390,
	0.0000000000, 0.0000000000, -0.392843978,
	0.0000000000, 0.0019687318, -0.005406536,
	0.0000000000, 0.7536376089, 0.000000000,
	-0.0068635336, 0.0120919868, 0.000000000,
	-1.2354254091, 0.6686678529, 0.000781085,
	0.0000000000, 0.0000000000, -0.005741567,
	0.0000000000, 0.0000000000, -0.013918444,
	0.0000000000, -0.0028960022, -0.008078630,
	0.0000000000, 0.0103118232, 0.037861147,
	0.0000000000, 0.0138055740, 0.000000000,
	-0.0008639247, -0.0047337405, 0.000000000,
	0.0092341743, 0.0205330307, 0.000000000,
	-0.0070845613, -0.0197922662, -0.005549342,
	0.0000000000, 0.0000000000, 0.406178935,
	0.0000000000, 0.0000000000, 0.018878926,
	0.0000000000, 0.0000000000, -0.419508519,
	0.0000000000, 0.0037680434, -0.016107128,
	0.0000000000, -0.7659181639, 1.204427794,
	0.0000000000, -0.0212213951, -0.001285688,
	0.0000000000, 0.7833715156, 0.000000000,
	-0.0246696555, 0.0000000000, -0.048223727,
	0.0388279407, 0.0000000000, -0.004262448,
	0.0000000000, 0.0790846570, -0.079420910,
	0.0000000000, 0.0132822485, 0.000000000,
	-0.0167973831, 0.0000000000, 0.000000000,
	0.0000000000, -0.0093383415, -0.011635432,
	0.0000000000, 0.0325317022, 0.000000000,
	0.0471128217, 0.0000000000, 0.008642167,
	-0.0664593336, -0.0303384405, 0.000000000,
	0.0000000000, -0.2515719173, 0.000000000,
	-0.0656383142, 0.0209441072, 0.000000000,
	1.3363896797, -0.0183893362, 0.000000000,
	0.0132897427, 0.0000000000, -0.034405130,
	0.0460133224, 0.0084317622, 0.000000000,
	0.0000000000, 0.1580348264, 0.000000000,
	-0.0396942406, 0.0406791300, 0.000000000,
	-0.6364590649, 0.0000000000, 0.012304309,
	-0.0464964922, 0.0000000000, -0.073662100,
	0.4496452159, 0.0000000000, -0.000209527,
	0.0000000000, 0.0039967903, 0.000000000,
	0.0000000000, 0.0000000000, 0.000646235,
	-0.0034051366, 0.0000000000, 0.001521557,
	0.0000000000, 0.0000000000, -0.004692869,
	0.0000000000, -0.0008036637, 0.000000000,
	0.0000000000, -0.0107033626, 0.000000000,
	0.0000000000, 0.0000000000, -0.002337854,
	0.0005581193, 0.0000000000, 0.006701593,
	0.0036573139, 0.0000000000, -0.002762281,
	0.0000000000, 0.0034502614, 0.000000000,
	-0.0113295274, 0.0061691646, 0.003646296,
	0.0000000000, 0.0000000000, 0.011958508,
	0.0000000000, -0.0049964428, 0.019105633,
	0.0000000000, -0.0007005211, 0.000000000,
	0.0120583878, -0.0115821314, 0.000000000,
	-0.0216924367, 0.0491315920, -0.002561508,
	0.0000000000, 0.0000000000, -0.002686464,
	0.0000000000, 0.0000000000, 0.004899140,
	0.0000000000, 0.0008639247, -0.004733740,
	0.0000000000, -0.0092341743, 0.020533030,
	0.0000000000, 0.0012856882, 0.002151334,
	0.0000000000, -0.0058027293, 0.000000000,
	0.0000000000, 0.0052554155, 0.000282545,
	0.0000000000, 0.0000000000, -0.005929873,
	0.0000000000, 0.0117230351, 0.000000000,
	0.0000000000, -0.0284442221, 0.000000000,
	0.0000000000, -0.0171561062, 0.000000000,
	-0.0000000000, 0.0312433423, -0.001216940,
	0.0000000000, 0.0000000000, 0.002823664,
	0.0000000000, 0.0000000000, -0.002704991,
	0.0000000000, -0.0029601529, 0.002970429,
	0.0000000000, 0.0083605935, -0.006708966,
	0.0000000000, 0.0039929761, 0.000000000,
	-0.0001995079, 0.0059503454, 0.000000000,
	0.0062913762, -0.0070436180, 0.000000000,
	-0.0011766655, -0.0081518477, 0.000000000,
	0.0115969342, 0.0000000000, -0.026027347,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0045901420, 0.0318379520, 0.000000000,
	-0.0037148194, 0.0000000000, 0.000000000,
	0.2866032276, 0.0000000000, 0.010761237,
	0.0000000000, 0.0000000000, -0.155763226,
	0.0000000000, 0.0000000000, 0.000000000,
	0.0019011469, -0.0011497730, 0.000000000,
	-0.0035302818, 0.0078497124, 0.000000000,
	0.0016774932, 0.0000000000, 0.000000000,
	-0.0129183402, 0.0000000000, 0.000000000,
	0.0003488328, 0.0000000000, -0.001216940,
	0.0000000000, 0.0000000000, 0.002823664,
	0.0000000000, 0.0000000000, 0.001098267,
	0.0000000000, 0.0000000000, 0.000000000,
	-0.0089079392, -0.0046608312, 0.000000000,
	-0.2793581265, 0.1371522760, 0.000000000,
	0.0108920142, -0.0027049915, 0.000000000,
	-0.0050997656, 0.0000000000, -0.002339355,
	0.0000000000, 0.0069498093, -0.018794488,
	0.0134098957, 0.0000000000, 0.000000000,
	0.2489517504, 0.0000000000, 0.023521935,
	-0.0486919075, 0.0000000000, -1.270050844,
	0.0000000000, 0.0143461100, 0.005308131,
	0.0000000000, 0.6978457422, 0.012226199,
	0.0012494780, 0.0000000000, 0.000000000,
	0.0067330993, 0.0000000000, 0.000000000,
	-0.0092628794, 0.0000000000, 0.000682256,
	-0.0030423654, 0.0000000000, 0.014931138,
	-0.0108989818, 0.0000000000, 0.007084561,
	0.0000000000, -0.0029601529, -0.002970429,
	0.0000000000, 0.0083605935, 0.006708966,
	0.0000000000, -0.0093934167, -0.004915202,
	-0.0095596081, 0.0000000000, 0.000000000,
	-0.2690947454, 0.0000000000, 0.000000000,
	0.0012803020, 0.0000000000, 0.000000000,
	0.2773740514, 0.0000000000, -0.021864836,
	0.0447844636, 0.0000000000, 1.231597770,
	-0.7012928704, 0.0000000000, -0.022697956,
	0.0039929761, 0.0000000000, -1.187034977,
	-0.0225899393, 0.0000000000, -0.003224564,
	0.0000000000, -0.0446795265, 0.028152172,
	0.0080830758, 0.0000000000, 0.000000000,
	-0.1404695206, 0.0000000000, 0.019329038,
	0.0257802037, 0.0000000000, 0.643959179,
	0.0000000000, 0.0221008213, 0.017639276,
	0.0000000000, 0.0329829706, -0.434392066,
	-0.0019761866, 0.0000000000, 0.000000000,
	0.0112848350, 0.0000000000, 0.000000000,
	-0.0015104417, 0.0000000000, 0.002231292,
	0.0024425376, 0.0000000000, -0.005773790,
	-0.0308404941, 0.0000000000, -0.019792266,
	0.0000000000, 0.0001995079, 0.005950345,
	0.0000000000, -0.0062913762, -0.007043618,
	0.0000000000, 0.0049152028, 0.009245120,
	0.0164830501, 0.0000000000, 0.000000000,
	0.1211016099, 0.0000000000, 0.000000000,
	-0.0077982066, 0.0000000000, 0.000000000,
	-0.1297864534, 0.0000000000, -0.018335766,
	0.0164567853, 0.0000000000, -0.657514428,
	-0.0242432979, 0.0000000000, 0.023334763,
	0.0011766655, 0.0000000000, 0.652515430,
	0.0000000000, 0.0223791973, -0.051741794,
	0.0000000000, -0.0524717981, 0.423796408,
	0.0000000000, 0.0234827537, -0.008151847,
	0.0000000000, 0.0066098471, -0.363902766,
}

func TestReadFortFile(t *testing.T) {
	t.Run("read from completed cartesian force constants", func(t *testing.T) {

		path, err := filepath.Abs("../testfiles/h2co/4th/fort.15")
		if err != nil {
			t.Fatal(err)
		}

		got := models.Chromosome(readFortFile(path))

		want := fort15

		if diff := deep.Equal(got, want); diff != nil {
			t.Error(diff)
		}
	})
}

func TestDNA(t *testing.T) {
	t.Run("get h2co dna from testfiles", func(t *testing.T) {

		dirPath, err := filepath.Abs(PathToH2CODir)
		if err != nil {
			t.Fatal(err)
		}

		got := ReadDNA(dirPath)

		if diff := deep.Equal(got[0], fort15); diff != nil {
			t.Error(diff)
		}

	})
}

func TestMockB3LYP(t *testing.T) {
	t.Run("is the difference within the parameters?", func(t *testing.T) {
		newDNA := MockB3LYP(PathToH2CODir)

		t.Run("compare with fort.15", func(t *testing.T) {
			assertBounds(t, 1.0, 0.90, newDNA[0], fort15)
		})

		t.Run("compare with fort.30", func(t *testing.T) {
			assertBounds(t, 3.0, 0.90, newDNA[1], fort30)
		})

		t.Run("compare with fort.40", func(t *testing.T) {
			assertBounds(t, 10.0, 0.90, newDNA[2], fort40)
		})
	})
}

func assertBounds(t *testing.T, upper float64, diffMax float64, new, old []float64) {
	for i, v := range new {
		diff := math.Abs(v - old[i])
		if diff > diffMax {
			t.Errorf("outside of testing parameters at index %d with val %v. Expected %v. The diff is %v\n", i, v, old[i], diff)
		}

		if math.Abs(v) > upper {
			t.Errorf("outside of testing bounds at index %d with val %v. Expected %v\n", i, v, old[i])
		}
	}
}