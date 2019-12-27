package schema

// THIS FILE WAS COPIED BY HAND FROM https://raw.githubusercontent.com/whosonfirst/whosonfirst-json-schema/master/schema/docs/wof-properties.json
// EVENTUALLY IT WILL BE CLONE (BY ROBOT) FROM THE SAME SOURCE. YOU SHOULD NOT UPDATE THIS FILE BY HAND.
// (20191227/thisisaaronland)

const PROPERTIES string = `{
  "$schema": "http://json-schema.org/draft-06/schema#",
  "$id": "wof-properties.json", 
  "definitions": {
    "properties": {
      "description": "The properties that can exist in a WOF document",
      "type": "object",
      "properties": {
        "wof:abbreviation": {
          "type": "string"
        },
        "wof:association": {
          "type": "string"
        },
        "wof:belongs": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:belongs_to": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:belongsto": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:brand_id": {
          "type": "string"
        },
        "wof:breaches": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:capital": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:capital_of": {
          "type": "string"
        },
        "wof:categories": {
          "type": "string"
        },
        "wof:category": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "wof:children": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:concordances": {
          "type": "object",
          "properties": {
            "gp:id": {
              "type": "integer"
            },
            "wd:id": {
              "type": "string"
            },
            "wk:page": {
              "type": "string"
            },
            "qs:id": {
              "type": "integer"
            },
            "loc:id": {
              "type": "string"
            },
            "qs_pg:id": {
              "type": "integer"
            },
            "gn:id": {
              "type": "integer"
            }
          }
        },
        "wof:concordances_alt": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "wof:concordances_sources": {
          "type": "string"
        },
        "wof:constituency": {
          "type": "string"
        },
        "wof:controlled": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "wof:coterminous": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:country": {
          "type": "string",
          "pattern": "^$|^[A-Za-z]{2}$|-99|-1"
        },
        "wof:country_alpha3": {
          "type": "string"
        },
        "wof:created": {
          "type": [
            "integer",
            "null"
          ]
        },
        "wof:fullname": {
          "type": "string"
        },
        "wof:geomhash": {
          "type": "string"
        },
        "wof:hierarchy": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "continent_id": {
                "type": "integer"
              },
              "locality_id": {
                "type": "integer"
              },
              "country_id": {
                "type": "integer"
              },
              "region_id": {
                "type": "integer"
              },
              "macroregion_id": {
                "type": "integer"
              }
            }
          }
        },
        "wof:id": {
          "type": "integer"
        },
        "wof:is_current": {
          "type": "string"
        },
        "wof:label": {
          "type": "string"
        },
        "wof:lang": {
          "type": "array",
          "items": {
            "pattern": "^[a-z]{3}$",
            "type": "string"
          }
        },
        "wof:lang_x_official": {
          "type": "array",
          "items": {
            "pattern": "^[a-z]{3}$",
            "type": "string"
          }
        },
        "wof:lang_x_spoken": {
          "type": "array",
          "items": {
            "pattern": "^[a-z]{3}$",
            "type": "string"
          }
        },
        "wof:lastmodified": {
          "type": "integer"
        },
        "wof:megacity": {
          "type": "integer"
        },
        "wof:name": {
          "type": [
            "string",
            "null"
          ]
        },
        "wof:parent_id": {
          "type": "integer"
        },
        "wof:phone": {
          "type": "string"
        },
        "wof:placetype": {
          "type": "string"
        },
        "wof:placetype_alt": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "wof:placetype_id": {
          "type": "string"
        },
        "wof:placetype_local": {
          "type": "string"
        },
        "wof:placetype_names": {
          "type": "string"
        },
        "wof:population": {
          "type": "integer"
        },
        "wof:population_rank": {
          "type": "integer"
        },
        "wof:repo": {
          "type": "string",
          "pattern": "^whosonfirst-.*$"
        },
        "wof:scale": {
          "type": "integer"
        },
        "wof:shortcode": {
          "type": "string"
        },
        "wof:statistical_gore": {
          "type": "integer"
        },
        "wof:subdivision": {
          "type": "string"
        },
        "wof:superseded": {
          "type": "string"
        },
        "wof:superseded_by": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:supersedes": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "wof:tags": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "wof:website": {
          "type": "string"
        },
        "abbreviation:eng_x_preferred": {
          "type": "array"
        },
        "abrv:eng_x_preferred": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "acgov:api": {
          "type": "string"
        },
        "acgov:type": {
          "type": "string"
        },
        "acme:elev": {
          "type": "string"
        },
        "acme:site_id": {
          "type": "string"
        },
        "addr:city": {
          "type": "string"
        },
        "addr:conscriptionnumber": {
          "type": "string"
        },
        "addr:country": {
          "type": "string"
        },
        "addr:district": {
          "type": "string"
        },
        "addr:email": {
          "type": "string"
        },
        "addr:facebook": {
          "type": "string"
        },
        "addr:flats": {
          "type": "string"
        },
        "addr:full": {
          "type": "string"
        },
        "addr:github": {
          "type": "string"
        },
        "addr:hamlet": {
          "type": "string"
        },
        "addr:housename": {
          "type": "string"
        },
        "addr:housenumber": {
          "type": "string"
        },
        "addr:instagram": {
          "type": "string"
        },
        "addr:intersection": {
          "type": "string"
        },
        "addr:notes": {
          "type": "string"
        },
        "addr:opentable": {
          "type": "string"
        },
        "addr:phon": {
          "type": "string"
        },
        "addr:phone": {
          "type": "string"
        },
        "addr:place": {
          "type": "string"
        },
        "addr:postal": {
          "type": "string"
        },
        "addr:postalcode": {
          "type": "string"
        },
        "addr:postcode": {
          "type": "string"
        },
        "addr:province": {
          "type": "string"
        },
        "addr:state": {
          "type": "string"
        },
        "addr:street": {
          "type": "string"
        },
        "addr:subdistrict": {
          "type": "string"
        },
        "addr:suburb": {
          "type": "string"
        },
        "addr:twitter": {
          "type": "string"
        },
        "addr:url": {
          "type": "string"
        },
        "addr:website": {
          "type": "string"
        },
        "addr:yelp": {
          "type": "string"
        },
        "addr:youtube": {
          "type": "string"
        },
        "amsgis:_categories": {
          "type": "string"
        },
        "amsgis:_id": {
          "type": "string"
        },
        "amsgis:categorie": {
          "type": "string"
        },
        "amsgis:CODE": {
          "type": "string"
        },
        "amsgis:display": {
          "type": "string"
        },
        "amsgis:DOCDATUM": {
          "type": "string"
        },
        "amsgis:DOCNR": {
          "type": "string"
        },
        "amsgis:externe_id": {
          "type": "string"
        },
        "amsgis:INGSDATUM": {
          "type": "string"
        },
        "amsgis:NAAM": {
          "type": "string"
        },
        "amsgis:online_tijdsaspect": {
          "type": "string"
        },
        "amsgis:titel": {
          "type": "string"
        },
        "amsgis:titel_key": {
          "type": "string"
        },
        "amsgis:type": {
          "type": "string"
        },
        "amsgis:type_2": {
          "type": "string"
        },
        "amsgis:uri": {
          "type": "string"
        },
        "amsgis:VOLLCODE": {
          "type": "string"
        },
        "atgov:bez_name": {
          "type": "string"
        },
        "atgov:bez_nr": {
          "type": "string"
        },
        "atgov:gem_name": {
          "type": "string"
        },
        "atgov:gkz": {
          "type": "string"
        },
        "atldpcd:FID_TXT": {
          "type": "string"
        },
        "atldpcd:GLOBALID": {
          "type": "string"
        },
        "atldpcd:NAME": {
          "type": "string"
        },
        "atldpcd:NPU": {
          "type": "string"
        },
        "atldpcd:OBJECTID": {
          "type": "string"
        },
        "atldpcd:OLD_NAME": {
          "type": "string"
        },
        "ausstat:POA_CODE": {
          "type": "string"
        },
        "ausstat:POA_NAME": {
          "type": "string"
        },
        "ausstat:SQKM": {
          "type": "string"
        },
        "austriaod:bez_name": {
          "type": "string"
        },
        "austriaod:bez_nr": {
          "type": "integer"
        },
        "austriaod:gem_name": {
          "type": "string"
        },
        "austriaod:gem_nr": {
          "type": "integer"
        },
        "austriaod:land_name": {
          "type": "string"
        },
        "austriaod:land_nr": {
          "type": "integer"
        },
        "austriaod:objectid": {
          "type": "integer"
        },
        "azavea:LISTNAME": {
          "type": "string"
        },
        "azavea:MAPNAME": {
          "type": "string"
        },
        "azavea:NAME": {
          "type": "string"
        },
        "baltomoit:label": {
          "type": "string"
        },
        "baltomoit:nbrdesc": {
          "type": "string"
        },
        "begov:AREA": {
          "type": "string"
        },
        "begov:BEGIN_LIFE": {
          "type": "string"
        },
        "begov:datpublbs": {
          "type": "string"
        },
        "begov:END_LIFE": {
          "type": "string"
        },
        "begov:ID": {
          "type": "string"
        },
        "begov:INSPIRE_ID": {
          "type": "string"
        },
        "begov:lengte": {
          "type": "number"
        },
        "begov:MU_ID": {
          "type": "string"
        },
        "begov:MU_NAME_DU": {
          "type": "string"
        },
        "begov:MU_NAME_FR": {
          "type": "string"
        },
        "begov:MU_NAT_COD": {
          "type": "string"
        },
        "begov:naam": {
          "type": "string"
        },
        "begov:NAT_CODE": {
          "type": "string"
        },
        "begov:niscode": {
          "type": "string"
        },
        "begov:numac": {
          "type": [
            "string",
            "null"
          ]
        },
        "begov:oidn": {
          "type": "string"
        },
        "begov:oppervl": {
          "type": "number"
        },
        "begov:PZ_ID": {
          "type": "string"
        },
        "begov:PZ_NAME_DU": {
          "type": "string"
        },
        "begov:PZ_NAME_FR": {
          "type": "string"
        },
        "begov:PZ_NAT_COD": {
          "type": "string"
        },
        "begov:terrid": {
          "type": "string"
        },
        "begov:uidn": {
          "type": "string"
        },
        "begov:VERSIONID": {
          "type": "string"
        },
        "bowie:latitude": {
          "type": "string"
        },
        "bowie:longitude": {
          "type": "string"
        },
        "bra:Name": {
          "type": "string"
        },
        "bra:Neighborho": {
          "type": "string"
        },
        "bra:OBJECTID": {
          "type": "string"
        },
        "camgov:N_HOOD": {
          "type": "string"
        },
        "camgov:NAME": {
          "type": "string"
        },
        "camgov:Webpage": {
          "type": "string"
        },
        "can-abog:CITY_ID": {
          "type": "integer"
        },
        "can-abog:GEOCODE": {
          "type": "string"
        },
        "can-abog:GEONAME": {
          "type": "string"
        },
        "can-abog:HAMLET_ID": {
          "type": "string"
        },
        "can-abog:PID": {
          "type": "integer"
        },
        "can-bbygov:NEIGHBOURH": {
          "type": "string"
        },
        "can-bbygov:OBJECTID_1": {
          "type": "string"
        },
        "can-bbygov:PSA": {
          "type": "string"
        },
        "can-calcai:class": {
          "type": "string"
        },
        "can-calcai:class_code": {
          "type": "string"
        },
        "can-calcai:comm_code": {
          "type": "string"
        },
        "can-calcai:comm_structure": {
          "type": "string"
        },
        "can-calcai:name": {
          "type": "string"
        },
        "can-calcai:sector": {
          "type": "string"
        },
        "can-calcai:srg": {
          "type": "string"
        },
        "can-dnvgov:GLOBALID": {
          "type": "string"
        },
        "can-dnvgov:MET_INPUT": {
          "type": "string"
        },
        "can-dnvgov:MET_TECH": {
          "type": "string"
        },
        "can-dnvgov:MET_TECH_R": {
          "type": "string"
        },
        "can-dnvgov:NBDY_NAME": {
          "type": "string"
        },
        "can-dnvgov:NBDY_NAME_": {
          "type": "string"
        },
        "can-dnvgov:OBJECTID": {
          "type": "string"
        },
        "can-dnvgov:STATS_ID": {
          "type": "string"
        },
        "can-dnvgov:YEAR": {
          "type": "string"
        },
        "can-edmdsd:name": {
          "type": "string"
        },
        "can-edmdsd:number": {
          "type": "string"
        },
        "can-gatsudd:LSECSTATID": {
          "type": "string"
        },
        "can-gatsudd:MUN_MRC": {
          "type": "string"
        },
        "can-gatsudd:NO_COMM": {
          "type": "string"
        },
        "can-gatsudd:NO_COMM_1": {
          "type": "string"
        },
        "can-gatsudd:NOM_COMM": {
          "type": "string"
        },
        "can-gatsudd:NOM_HISTOR": {
          "type": "string"
        },
        "can-gatsudd:POP_2006": {
          "type": "string"
        },
        "can-gatsudd:SECTEUR": {
          "type": "string"
        },
        "can-gatsudd:SECTREG_2": {
          "type": "string"
        },
        "can-mntsmvt:no_arr": {
          "type": "string"
        },
        "can-mntsmvt:no_qr": {
          "type": "string"
        },
        "can-mntsmvt:nom_arr": {
          "type": "string"
        },
        "can-mntsmvt:nom_mun": {
          "type": "string"
        },
        "can-mntsmvt:nom_qr": {
          "type": "string"
        },
        "can-nwds:NEIGH_NAME": {
          "type": "string"
        },
        "can-nwds:NEIGHNUM": {
          "type": "string"
        },
        "can-ons:Name": {
          "type": "string"
        },
        "can-ons:Name2016": {
          "type": "string"
        },
        "can-ons:Name2016_F": {
          "type": "string"
        },
        "can-ons:Name2017": {
          "type": "string"
        },
        "can-ons:ONSID": {
          "type": "string"
        },
        "can-wpgppd:id": {
          "type": "string"
        },
        "can-wpgppd:name": {
          "type": "string"
        },
        "canvec-hydro:definit": {
          "type": "string"
        },
        "canvec-hydro:definit_en": {
          "type": "string"
        },
        "cbsnl:BU_CODE": {
          "type": "string"
        },
        "cbsnl:BU_NAAM": {
          "type": "string"
        },
        "cbsnl:GM_CODE": {
          "type": "string"
        },
        "cbsnl:GM_NAAM": {
          "type": "string"
        },
        "cbsnl:IND_WBI": {
          "type": "string"
        },
        "cbsnl:OAD": {
          "type": "string"
        },
        "cbsnl:STED": {
          "type": "string"
        },
        "cbsnl:WATER": {
          "type": "string"
        },
        "cbsnl:WK_CODE": {
          "type": "string"
        },
        "cbsnl:WK_NAAM": {
          "type": "string"
        },
        "chgov:os_uuid": {
          "type": "string"
        },
        "chgov:uuid": {
          "type": "string"
        },
        "clustr:alpha": {
          "type": "string"
        },
        "clustr:area": {
          "type": "string"
        },
        "clustr:count": {
          "type": "string"
        },
        "clustr:density": {
          "type": "string"
        },
        "clustr:perimeter": {
          "type": "string"
        },
        "counts:concordances_total": {
          "type": "string"
        },
        "counts:languages_official": {
          "type": "string"
        },
        "counts:languages_spoken": {
          "type": "string"
        },
        "counts:languages_total": {
          "type": "string"
        },
        "counts:names_colloquial": {
          "type": "string"
        },
        "counts:names_languages": {
          "type": "string"
        },
        "counts:names_prefered": {
          "type": "string"
        },
        "counts:names_total": {
          "type": "string"
        },
        "counts:names_variant": {
          "type": "string"
        },
        "denvercpd:NBHD_ID": {
          "type": "string"
        },
        "denvercpd:NBHD_NAME": {
          "type": "string"
        },
        "ebc:bdyset_id": {
          "type": "string"
        },
        "ebc:ed_abbrev": {
          "type": "string"
        },
        "ebc:ed_id": {
          "type": "string"
        },
        "ebc:ed_name": {
          "type": "string"
        },
        "ebc:feat_area": {
          "type": "string"
        },
        "ebc:feat_perim": {
          "type": "string"
        },
        "ebc:gazette_dt": {
          "type": "string"
        },
        "ebc:objectid": {
          "type": "string"
        },
        "edtf:cessation": {
          "type": "string"
        },
        "edtf:deprecate": {
          "type": "string"
        },
        "edtf:deprecated": {
          "type": "string"
        },
        "edtf:inception": {
          "type": "string"
        },
        "edtf:superseded": {
          "type": "string"
        },
        "esp-aytomad:CODBAR": {
          "type": "string"
        },
        "esp-aytomad:CODBARRIO": {
          "type": "string"
        },
        "esp-aytomad:CODDISTRIT": {
          "type": "string"
        },
        "esp-aytomad:NOMBRE": {
          "type": "string"
        },
        "esp-aytomad:NOMDIS": {
          "type": "string"
        },
        "esp-aytomad:OBJECTID": {
          "type": "string"
        },
        "esp-cartobcn:C_Barri": {
          "type": "string"
        },
        "esp-cartobcn:C_Distri": {
          "type": "string"
        },
        "esp-cartobcn:N_Barri": {
          "type": "string"
        },
        "esp-cartobcn:N_Distri": {
          "type": "string"
        },
        "esp-cartobcn:WEB_1": {
          "type": "string"
        },
        "esp-cartobcn:WEB_4": {
          "type": "string"
        },
        "figov:ajo_pvm": {
          "type": "string"
        },
        "figov:Aluejako": {
          "type": "string"
        },
        "figov:eng_type": {
          "type": "string"
        },
        "figov:fin_type": {
          "type": "string"
        },
        "figov:gml_id": {
          "type": "string"
        },
        "figov:Kunta": {
          "type": "string"
        },
        "figov:local_id": {
          "type": "string"
        },
        "figov:national_code": {
          "type": "string"
        },
        "figov:nimi": {
          "type": "string"
        },
        "figov:nimi_se": {
          "type": "string"
        },
        "figov:swe_type": {
          "type": "string"
        },
        "figov:Tunnus": {
          "type": "string"
        },
        "fra-odp:c_ar": {
          "type": "string"
        },
        "fra-odp:c_qu": {
          "type": "string"
        },
        "fra-odp:c_quinsee": {
          "type": "string"
        },
        "fra-odp:l_qu": {
          "type": "string"
        },
        "frgov:_COL6": {
          "type": "string"
        },
        "frgov:DEP": {
          "type": "string"
        },
        "frgov:ID": {
          "type": "string"
        },
        "frgov:LIB": {
          "type": "string"
        },
        "frgov:POP2010": {
          "type": "string"
        },
        "frgov:SURF": {
          "type": "string"
        },
        "fsgov:ajo_pvm": {
          "type": "string"
        },
        "fsgov:Aluejako": {
          "type": "string"
        },
        "fsgov:Kunta": {
          "type": "string"
        },
        "fsgov:nimi": {
          "type": "string"
        },
        "fsgov:nimi_se": {
          "type": "string"
        },
        "fsgov:Tunnus": {
          "type": "string"
        },
        "gbr-datalondon:GSS_CODE": {
          "type": "string"
        },
        "gbr-datalondon:ONS_INNER": {
          "type": "string"
        },
        "geom:area": {
          "type": "number"
        },
        "geom:area_square_m": {
          "type": "number"
        },
        "geom:bbox": {
          "type": "string",
          "pattern": "^([-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?),([-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?),([-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?),([-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?)$"
        },
        "geom:hash": {
          "type": "string"
        },
        "geom:latitude": {
          "type": "number"
        },
        "geom:longitude": {
          "type": "number"
        },
        "geom:src": {
          "type": "string"
        },
        "geom:type": {
          "type": "string"
        },
        "geonames:id": {
          "type": "string"
        },
        "gn:accuracy": {
          "type": "string"
        },
        "gn:adm1_code": {
          "type": "string"
        },
        "gn:adm1_name": {
          "type": "string"
        },
        "gn:adm2_code": {
          "type": "string"
        },
        "gn:adm2_name": {
          "type": "string"
        },
        "gn:adm3_code": {
          "type": "string"
        },
        "gn:adm3_name": {
          "type": "string"
        },
        "gn:country": {
          "type": "string"
        },
        "gn:elevation": {
          "type": "integer"
        },
        "gn:fcode": {
          "type": "string"
        },
        "gn:gn_country": {
          "type": [
            "string",
            "null"
          ]
        },
        "gn:gn_fcode": {
          "type": "string"
        },
        "gn:gn_pop": {
          "type": "integer"
        },
        "gn:id": {
          "type": "string"
        },
        "gn:latitude": {
          "type": "number"
        },
        "gn:longitude": {
          "type": "number"
        },
        "gn:name": {
          "type": "string"
        },
        "gn:pop": {
          "type": "integer"
        },
        "gn:population": {
          "type": "integer"
        },
        "goem:longitude": {
          "type": "number"
        },
        "gp:adm0": {
          "type": "integer"
        },
        "gp:id": {
          "type": "integer"
        },
        "gp:parent_id": {
          "type": "integer"
        },
        "gp:source": {
          "type": "string"
        },
        "hkigis:ajo_pvm": {
          "type": "string"
        },
        "hkigis:aluejako": {
          "type": "string"
        },
        "hkigis:ALUETASO": {
          "type": "string"
        },
        "hkigis:ID": {
          "type": "string"
        },
        "hkigis:ID1": {
          "type": "string"
        },
        "hkigis:K_NIMI_SE": {
          "type": "string"
        },
        "hkigis:KOKOTUN": {
          "type": "string"
        },
        "hkigis:KOKOTUNNUS": {
          "type": "string"
        },
        "hkigis:KUNTA": {
          "type": "string"
        },
        "hkigis:KUNTA_NIMI": {
          "type": "string"
        },
        "hkigis:Mtryhm": {
          "type": "string"
        },
        "hkigis:NIMI": {
          "type": "string"
        },
        "hkigis:NIMI_ISO": {
          "type": "string"
        },
        "hkigis:NIMI_SE": {
          "type": "string"
        },
        "hkigis:PIEN": {
          "type": "string"
        },
        "hkigis:SUUR": {
          "type": "string"
        },
        "hkigis:SUUR_N_FI": {
          "type": "string"
        },
        "hkigis:SUUR_N_SE": {
          "type": "string"
        },
        "hkigis:SUURP_TN": {
          "type": "string"
        },
        "hkigis:TILA": {
          "type": "string"
        },
        "hkigis:TUNNUS": {
          "type": "string"
        },
        "intersection:latitude": {
          "type": "number"
        },
        "intersection:longitude": {
          "type": "number"
        },
        "iso:country": {
          "type": "string",
          "pattern": "^$|^[A-Za-z]{2}$|-99|-1"
        },
        "iso:parent": {
          "type": "string"
        },
        "iso:subdivision": {
          "type": "string"
        },
        "itu:country_code": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "itu:region": {
          "type": "string"
        },
        "kuogov:ID": {
          "type": "string"
        },
        "kuogov:NIMI": {
          "type": "string"
        },
        "label:{lang}_x_colloquial": {
          "type": "array"
        },
        "label:{lang}_x_colloquial_abbreviation": {
          "type": "array"
        },
        "label:{lang}_x_historical": {
          "type": "array"
        },
        "label:{lang}_x_preferred": {
          "type": "array"
        },
        "label:{lang}_x_preferred_abbreviation": {
          "type": "array"
        },
        "label:{lang}_x_preferred_disambiguation": {
          "type": "array"
        },
        "label:{lang}_x_preferred_placetype": {
          "type": "array"
        },
        "label:{lang}_x_preferred_shortcode": {
          "type": "array"
        },
        "label:{lang}_x_unknown": {
          "type": "array"
        },
        "label:{lang}_x_variant": {
          "type": "array"
        },
        "label:{lang}_x_variant_abbreviation": {
          "type": "array"
        },
        "label:{lang}_x_variant_disambiguation": {
          "type": "array"
        },
        "label:{lang}_x_variant_placetype": {
          "type": "array"
        },
        "label:{lang}_x_variant_shortcode": {
          "type": "array"
        },
        "lacity:CERTIFIED": {
          "type": "string"
        },
        "lacity:DWEBSITE": {
          "type": "string"
        },
        "lacity:NAME": {
          "type": "string"
        },
        "lacity:NC_ID": {
          "type": "string"
        },
        "lacity:NSA": {
          "type": "string"
        },
        "lacity:OBJECTID": {
          "type": "string"
        },
        "lacity:WADDRESS": {
          "type": "string"
        },
        "lbl:bbox": {
          "type": "string"
        },
        "lbl:latitude": {
          "type": "number"
        },
        "lbl:longitude": {
          "type": "number"
        },
        "lflt:label_text": {
          "type": "string"
        },
        "local:latitude": {
          "type": "number"
        },
        "local:longitude": {
          "type": "number"
        },
        "meso:adm0": {
          "type": "string"
        },
        "meso:adm1": {
          "type": "string"
        },
        "meso:adm1_alt": {
          "type": "string"
        },
        "meso:adm1_alt2": {
          "type": "string"
        },
        "meso:adm1_en": {
          "type": "string"
        },
        "meso:adm1_loc": {
          "type": "string"
        },
        "meso:adm1_loc2": {
          "type": "string"
        },
        "meso:adm1_name": {
          "type": "string"
        },
        "meso:adm2_alt": {
          "type": "string"
        },
        "meso:adm2_gaul": {
          "type": "string"
        },
        "meso:adm2_loc": {
          "type": "string"
        },
        "meso:adm2_name": {
          "type": "string"
        },
        "meso:admin1": {
          "type": "string"
        },
        "meso:admin1_en": {
          "type": "string"
        },
        "meso:admin1_loc": {
          "type": "string"
        },
        "meso:admin1name": {
          "type": "string"
        },
        "meso:admin1r": {
          "type": "string"
        },
        "meso:admin1r_en": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:admin2name": {
          "type": "string"
        },
        "meso:admin_1": {
          "type": "string"
        },
        "meso:admin_1_en": {
          "type": "string"
        },
        "meso:admin_1_lo": {
          "type": "string"
        },
        "meso:admin_1r": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:admin_2": {
          "type": "string"
        },
        "meso:aotm": {
          "type": "string"
        },
        "meso:c_desc": {
          "type": "string"
        },
        "meso:c_name": {
          "type": "string"
        },
        "meso:county_id": {
          "type": "string"
        },
        "meso:delete_fla": {
          "type": "string"
        },
        "meso:diss_me": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:dissme": {
          "type": "string"
        },
        "meso:gadm_adm1": {
          "type": "string"
        },
        "meso:gadm_adm2": {
          "type": "string"
        },
        "meso:gadm_adm3": {
          "type": "string"
        },
        "meso:gadm_admi1": {
          "type": "string"
        },
        "meso:gadm_admin": {
          "type": "string"
        },
        "meso:gadm_admn2": {
          "type": "string"
        },
        "meso:gadm_alt": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:gadm_cyril": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:gadm_hasc": {
          "type": "string"
        },
        "meso:gadm_hasc2": {
          "type": "string"
        },
        "meso:gadm_hasc_": {
          "type": "string"
        },
        "meso:gadm_loc": {
          "type": "string"
        },
        "meso:gadm_loc_a": {
          "type": "string"
        },
        "meso:gadm_name": {
          "type": "string"
        },
        "meso:gadm_name_": {
          "type": "string"
        },
        "meso:gaul_adm1": {
          "type": "string"
        },
        "meso:gaul_adm2": {
          "type": "string"
        },
        "meso:gaul_adm2_": {
          "type": "string"
        },
        "meso:gaul_adm3": {
          "type": "string"
        },
        "meso:gaul_admin": {
          "type": "string"
        },
        "meso:gaul_name": {
          "type": "string"
        },
        "meso:hasc_1": {
          "type": "string"
        },
        "meso:hasc_id": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:id": {
          "type": "string"
        },
        "meso:iscgm_adm2": {
          "type": "string"
        },
        "meso:iscgm_admi": {
          "type": "string"
        },
        "meso:loc_id_alt": {
          "type": "string"
        },
        "meso:local_id": {
          "type": "string"
        },
        "meso:local_id_a": {
          "type": "string"
        },
        "meso:mps_x": {
          "type": "number"
        },
        "meso:mps_y": {
          "type": "number"
        },
        "meso:nam_en_alt": {
          "type": "string"
        },
        "meso:name_alt": {
          "type": "string"
        },
        "meso:name_alt1": {
          "type": "string"
        },
        "meso:name_alt2": {
          "type": "string"
        },
        "meso:name_alt3": {
          "type": "string"
        },
        "meso:name_en": {
          "type": "string"
        },
        "meso:name_en2": {
          "type": "string"
        },
        "meso:name_lat": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:name_loc": {
          "type": "string"
        },
        "meso:name_loc2": {
          "type": "string"
        },
        "meso:name_loc3": {
          "type": "string"
        },
        "meso:name_loc_1": {
          "type": "string"
        },
        "meso:name_loc_a": {
          "type": "string"
        },
        "meso:name_local": {
          "type": "string"
        },
        "meso:ne_adm1_en": {
          "type": "string"
        },
        "meso:ne_adm1_lo": {
          "type": "string"
        },
        "meso:ne_code_ha": {
          "type": "string"
        },
        "meso:objectid": {
          "type": "integer"
        },
        "meso:ocha_adm2": {
          "type": "string"
        },
        "meso:ocha_admin": {
          "type": "string"
        },
        "meso:okato_code": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:okato_name": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:oktmo_code": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:oktmo_name": {
          "type": [
            "string",
            "null"
          ]
        },
        "meso:osm_adm2": {
          "type": "string"
        },
        "meso:patch": {
          "type": "integer"
        },
        "meso:pop": {
          "type": "integer"
        },
        "meso:pop_year": {
          "type": "string"
        },
        "meso:region": {
          "type": "string"
        },
        "meso:source": {
          "type": "string"
        },
        "meso:type_alt": {
          "type": "string"
        },
        "meso:type_en": {
          "type": "string"
        },
        "meso:type_loc": {
          "type": "string"
        },
        "meso:type_local": {
          "type": "string"
        },
        "misc:": {
          "type": "string"
        },
        "misc:gn_adm0_cc": {
          "type": "string"
        },
        "misc:gn_fcode": {
          "type": "string"
        },
        "misc:gn_id": {
          "type": "integer"
        },
        "misc:gn_local": {
          "type": "integer"
        },
        "misc:gn_nam_loc": {
          "type": "string"
        },
        "misc:gn_namadm1": {
          "type": "string"
        },
        "misc:gn_name": {
          "type": "string"
        },
        "misc:local_max": {
          "type": "integer"
        },
        "misc:local_sum": {
          "type": "integer"
        },
        "misc:localhoods": {
          "type": "integer"
        },
        "misc:name": {
          "type": "string"
        },
        "misc:name_adm0": {
          "type": "string"
        },
        "misc:name_adm1": {
          "type": "string"
        },
        "misc:name_adm2": {
          "type": "string"
        },
        "misc:name_en": {
          "type": "string"
        },
        "misc:name_lau": {
          "type": "string"
        },
        "misc:name_local": {
          "type": "string"
        },
        "misc:photo_max": {
          "type": "integer"
        },
        "misc:photo_sum": {
          "type": "integer"
        },
        "misc:placetype": {
          "type": "string"
        },
        "misc:qs_a0": {
          "type": "string"
        },
        "misc:qs_a0_lc": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_a1": {
          "type": "string"
        },
        "misc:qs_a1_alt": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_a1_lc": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_a1r": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_a1r_alt": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_a1r_lc": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_adm0": {
          "type": "string"
        },
        "misc:qs_adm0_a3": {
          "type": "string"
        },
        "misc:qs_gn_id": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_id": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_iso_cc": {
          "type": "string"
        },
        "misc:qs_level": {
          "type": "string"
        },
        "misc:qs_pop": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_scale": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:qs_source": {
          "type": "string"
        },
        "misc:qs_type": {
          "type": "string"
        },
        "misc:qs_woe_id": {
          "type": [
            "string",
            "null"
          ]
        },
        "misc:quad_count": {
          "type": "integer"
        },
        "misc:woe_adm0": {
          "type": "integer"
        },
        "misc:woe_adm1": {
          "type": "integer"
        },
        "misc:woe_adm2": {
          "type": "integer"
        },
        "misc:woe_funk": {
          "type": "string"
        },
        "misc:woe_lau": {
          "type": "integer"
        },
        "misc:woe_local": {
          "type": "integer"
        },
        "misc:woe_ver": {
          "type": "string"
        },
        "misc:woeid": {
          "type": "integer"
        },
        "mps:latitude": {
          "type": "number"
        },
        "mps:longitude": {
          "type": "number"
        },
        "mz:berth": {
          "type": "string"
        },
        "mz:categories": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "mz:hierarchy_label": {
          "type": "integer"
        },
        "mz:hours": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "mz:is_approximate": {
          "type": "integer"
        },
        "mz:is_clustr": {
          "type": "integer"
        },
        "mz:is_current": {
          "type": "integer"
        },
        "mz:is_funky": {
          "type": "integer"
        },
        "mz:is_hard_boundary": {
          "type": "integer"
        },
        "mz:is_landuse_aoi": {
          "type": "integer"
        },
        "mz:is_official": {
          "type": "integer"
        },
        "mz:is_retail": {
          "type": "integer"
        },
        "mz:max_zoom": {
          "type": "number"
        },
        "mz:min_zoom": {
          "type": "number"
        },
        "mz:note": {
          "type": "string"
        },
        "mz:phone": {
          "type": "string"
        },
        "mz:remarks": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "mz:tier_locality": {
          "type": "integer"
        },
        "mz:tier_metro": {
          "type": "integer"
        },
        "nav:latitude": {
          "type": "number"
        },
        "nav:longitude": {
          "type": "number"
        },
        "ne:abbrev": {
          "type": "string"
        },
        "ne:abbrev_len": {
          "type": "integer"
        },
        "ne:ADM0_A3": {
          "type": "string",
          "pattern": "^[A-Z]{3}$"
        },
        "ne:adm0_a3_is": {
          "type": "string"
        },
        "ne:adm0_a3_un": {
          "type": "integer"
        },
        "ne:adm0_a3_us": {
          "type": "string"
        },
        "ne:adm0_a3_wb": {
          "type": "integer"
        },
        "ne:adm0_dif": {
          "type": "integer"
        },
        "ne:adm0_label": {
          "type": "string"
        },
        "ne:adm0_sr": {
          "type": "string"
        },
        "ne:ADM0CAP": {
          "type": "number"
        },
        "ne:ADM0NAME": {
          "type": "string"
        },
        "ne:adm1_cod_1": {
          "type": "string"
        },
        "ne:adm1_code": {
          "type": "string"
        },
        "ne:ADM1NAME": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:admin": {
          "type": "string"
        },
        "ne:ADMIN1_COD": {
          "type": "integer"
        },
        "ne:area_sqkm": {
          "type": "string"
        },
        "ne:brk_a3": {
          "type": "string"
        },
        "ne:brk_diff": {
          "type": "integer"
        },
        "ne:brk_group": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:brk_name": {
          "type": "string"
        },
        "ne:CAPALT": {
          "type": "integer"
        },
        "ne:CAPIN": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:CHANGED": {
          "type": "integer"
        },
        "ne:check_me": {
          "type": "string"
        },
        "ne:CHECKME": {
          "type": "integer"
        },
        "ne:CITYALT": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:code_hasc": {
          "type": "string"
        },
        "ne:code_local": {
          "type": "string"
        },
        "ne:COMPARE": {
          "type": "integer"
        },
        "ne:continent": {
          "type": "string"
        },
        "ne:datarank": {
          "type": "string"
        },
        "ne:DIFFASCII": {
          "type": "integer"
        },
        "ne:DIFFNOTE": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:diss_me": {
          "type": "string"
        },
        "ne:economy": {
          "type": "string"
        },
        "ne:ELEVATION": {
          "type": "number"
        },
        "ne:FEATURE_CL": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:FEATURE_CO": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:FEATURECLA": {
          "type": "string"
        },
        "ne:fips": {
          "type": "string"
        },
        "ne:fips_10": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:fips_10_": {
          "type": "string"
        },
        "ne:fips_alt": {
          "type": "string"
        },
        "ne:formal_en": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:formal_fr": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:gadm_level": {
          "type": "string"
        },
        "ne:gdp_md_est": {
          "type": "integer"
        },
        "ne:gdp_year": {
          "type": "integer"
        },
        "ne:GEONAMEID": {
          "type": "integer"
        },
        "ne:GEONAMESNO": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:geonunit": {
          "type": "string"
        },
        "ne:geou_dif": {
          "type": "integer"
        },
        "ne:geounit": {
          "type": "string"
        },
        "ne:gn_a1_code": {
          "type": "string"
        },
        "ne:GN_ASCII": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:gn_id": {
          "type": "string"
        },
        "ne:gn_level": {
          "type": "string"
        },
        "ne:gn_name": {
          "type": "string"
        },
        "ne:GN_POP": {
          "type": "integer"
        },
        "ne:gn_region": {
          "type": "string"
        },
        "ne:gns_adm1": {
          "type": "string"
        },
        "ne:gns_id": {
          "type": "string"
        },
        "ne:gns_lang": {
          "type": "string"
        },
        "ne:gns_level": {
          "type": "string"
        },
        "ne:gns_name": {
          "type": "string"
        },
        "ne:gns_region": {
          "type": "string"
        },
        "ne:GTOPO30": {
          "type": "number"
        },
        "ne:gu_a3": {
          "type": "string"
        },
        "ne:hasc_maybe": {
          "type": "string"
        },
        "ne:homepart": {
          "type": "integer"
        },
        "ne:income_grp": {
          "type": "string"
        },
        "ne:iso_3166_2": {
          "type": "string"
        },
        "ne:ISO_A2": {
          "type": "string",
          "pattern": "^$|^[A-Za-z]{2}$|-99|-1"
        },
        "ne:iso_a3": {
          "type": "string"
        },
        "ne:iso_n3": {
          "type": "integer"
        },
        "ne:LABELRANK": {
          "type": "integer"
        },
        "ne:lastcensus": {
          "type": "integer"
        },
        "ne:LATITUDE": {
          "type": "number"
        },
        "ne:level": {
          "type": "integer"
        },
        "ne:long_len": {
          "type": "integer"
        },
        "ne:LONGITUDE": {
          "type": "number"
        },
        "ne:LS_MATCH": {
          "type": "integer"
        },
        "ne:LS_NAME": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:map_color": {
          "type": "number"
        },
        "ne:mapcolor13": {
          "type": "integer"
        },
        "ne:mapcolor7": {
          "type": "integer"
        },
        "ne:mapcolor8": {
          "type": "integer"
        },
        "ne:mapcolor9": {
          "type": "integer"
        },
        "ne:MAX_AREAKM": {
          "type": "number"
        },
        "ne:MAX_AREAMI": {
          "type": "number"
        },
        "ne:MAX_BBXMAX": {
          "type": "number"
        },
        "ne:MAX_BBXMIN": {
          "type": "number"
        },
        "ne:MAX_BBYMAX": {
          "type": "number"
        },
        "ne:MAX_BBYMIN": {
          "type": "number"
        },
        "ne:MAX_NATSCA": {
          "type": "number"
        },
        "ne:MAX_PERKM": {
          "type": "number"
        },
        "ne:MAX_PERMI": {
          "type": "number"
        },
        "ne:MAX_POP10": {
          "type": "integer"
        },
        "ne:MAX_POP20": {
          "type": "integer"
        },
        "ne:MAX_POP300": {
          "type": "integer"
        },
        "ne:MAX_POP310": {
          "type": "integer"
        },
        "ne:MAX_POP50": {
          "type": "integer"
        },
        "ne:MEAN_BBXC": {
          "type": "number"
        },
        "ne:MEAN_BBYC": {
          "type": "number"
        },
        "ne:MEGACITY": {
          "type": "integer"
        },
        "ne:MEGANAME": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:MIN_AREAKM": {
          "type": "number"
        },
        "ne:MIN_AREAMI": {
          "type": "number"
        },
        "ne:MIN_BBXMAX": {
          "type": "number"
        },
        "ne:MIN_BBXMIN": {
          "type": "number"
        },
        "ne:MIN_BBYMAX": {
          "type": "number"
        },
        "ne:MIN_BBYMIN": {
          "type": "number"
        },
        "ne:MIN_PERKM": {
          "type": "number"
        },
        "ne:MIN_PERMI": {
          "type": "number"
        },
        "ne:min_zoom": {
          "type": "number"
        },
        "ne:NAME": {
          "type": "string"
        },
        "ne:name_alt": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:name_forma": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:name_fr": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:name_len": {
          "type": "integer"
        },
        "ne:name_local": {
          "type": "string"
        },
        "ne:name_long": {
          "type": "string"
        },
        "ne:name_sort": {
          "type": "string"
        },
        "ne:NAMEALT": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:NAMEASCII": {
          "type": "string"
        },
        "ne:NAMEDIFF": {
          "type": "integer"
        },
        "ne:NAMEPAR": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:NATSCALE": {
          "type": "integer"
        },
        "ne:ne_10m_adm": {
          "type": "string"
        },
        "ne:NOTE": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:note_adm0": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:note_brk": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:oid_": {
          "type": "integer"
        },
        "ne:POP1950": {
          "type": "integer"
        },
        "ne:POP1955": {
          "type": "integer"
        },
        "ne:POP1960": {
          "type": "integer"
        },
        "ne:POP1965": {
          "type": "integer"
        },
        "ne:POP1970": {
          "type": "integer"
        },
        "ne:POP1975": {
          "type": "integer"
        },
        "ne:POP1980": {
          "type": "integer"
        },
        "ne:POP1985": {
          "type": "integer"
        },
        "ne:POP1990": {
          "type": "integer"
        },
        "ne:POP1995": {
          "type": "integer"
        },
        "ne:POP2000": {
          "type": "integer"
        },
        "ne:POP2005": {
          "type": "integer"
        },
        "ne:POP2010": {
          "type": "integer"
        },
        "ne:POP2015": {
          "type": "integer"
        },
        "ne:POP2020": {
          "type": "integer"
        },
        "ne:POP2025": {
          "type": "integer"
        },
        "ne:POP2050": {
          "type": "integer"
        },
        "ne:pop_est": {
          "type": "integer"
        },
        "ne:POP_MAX": {
          "type": "integer"
        },
        "ne:POP_MIN": {
          "type": "integer"
        },
        "ne:POP_OTHER": {
          "type": "integer"
        },
        "ne:pop_year": {
          "type": "integer"
        },
        "ne:postal": {
          "type": "string"
        },
        "ne:provnum_ne": {
          "type": "string"
        },
        "ne:RANK_MAX": {
          "type": "integer"
        },
        "ne:RANK_MIN": {
          "type": "integer"
        },
        "ne:region": {
          "type": "string"
        },
        "ne:region_cod": {
          "type": "string"
        },
        "ne:region_sub": {
          "type": "string"
        },
        "ne:region_un": {
          "type": "string"
        },
        "ne:region_wb": {
          "type": "string"
        },
        "ne:sameascity": {
          "type": "string"
        },
        "ne:SCALERANK": {
          "type": "integer"
        },
        "ne:SOV0NAME": {
          "type": "string"
        },
        "ne:SOV_A3": {
          "type": "string",
          "pattern": "^$|^[A-Za-z1-9]{3}$|-99|-1"
        },
        "ne:sovereignt": {
          "type": "string"
        },
        "ne:su_a3": {
          "type": "string"
        },
        "ne:su_dif": {
          "type": "integer"
        },
        "ne:sub_code": {
          "type": "string"
        },
        "ne:subregion": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:subunit": {
          "type": "string"
        },
        "ne:terr_": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:TIMEZONE": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:tiny": {
          "type": "integer"
        },
        "ne:type": {
          "type": "string"
        },
        "ne:type_en": {
          "type": "string"
        },
        "ne:un_a3": {
          "type": "string"
        },
        "ne:UN_ADM0": {
          "type": [
            "string",
            "null"
          ]
        },
        "ne:UN_FID": {
          "type": "integer"
        },
        "ne:UN_LAT": {
          "type": "number"
        },
        "ne:UN_LONG": {
          "type": "number"
        },
        "ne:wb_a2": {
          "type": "string"
        },
        "ne:wb_a3": {
          "type": "string"
        },
        "ne:wikipedia": {
          "type": "string"
        },
        "ne:woe_id": {
          "type": "integer"
        },
        "ne:woe_id_eh": {
          "type": "integer"
        },
        "ne:woe_label": {
          "type": "string"
        },
        "ne:woe_name": {
          "type": "string"
        },
        "ne:woe_note": {
          "type": "string"
        },
        "ne:WORLDCITY": {
          "type": "integer"
        },
        "nolagis:gnocdc_lab": {
          "type": "string"
        },
        "nolagis:lup_lab": {
          "type": "string"
        },
        "nolagis:neigh_id": {
          "type": "string"
        },
        "oa:elevation_ft": {
          "type": "string"
        },
        "oa:type": {
          "type": "string"
        },
        "oakced:name": {
          "type": "string"
        },
        "os:admin_county_code": {
          "type": "string"
        },
        "os:admin_distict_code": {
          "type": "string"
        },
        "os:admin_ward_code": {
          "type": "string"
        },
        "os:country_code": {
          "type": "string"
        },
        "os:nhs_ha_code": {
          "type": "string"
        },
        "os:nhs_regional_ha_code": {
          "type": "string"
        },
        "os:positional_quality_indicator": {
          "type": "string"
        },
        "oulugov:AJ_KAUPU00": {
          "type": "string"
        },
        "oulugov:AJ_KAUPUNG": {
          "type": "string"
        },
        "oulugov:ID": {
          "type": "string"
        },
        "oulugov:TEKSTI": {
          "type": "string"
        },
        "out:cartodb_id": {
          "type": "string"
        },
        "out:image": {
          "type": "string"
        },
        "out:mafia_owned": {
          "type": "string"
        },
        "out:membership": {
          "type": "string"
        },
        "out:notes": {
          "type": "string"
        },
        "out:pub_desc": {
          "type": "string"
        },
        "out:source": {
          "type": "string"
        },
        "out:the_geom": {
          "type": "string"
        },
        "pedia:@id": {
          "type": "string"
        },
        "pedia:borough": {
          "type": "string"
        },
        "pedia:boroughCode": {
          "type": "string"
        },
        "pedia:neighborhood": {
          "type": "string"
        },
        "porbps:COALIT": {
          "type": "string"
        },
        "porbps:COMMPLAN": {
          "type": "string"
        },
        "porbps:HORZ_VERT": {
          "type": "string"
        },
        "porbps:ID": {
          "type": "string"
        },
        "porbps:MAPLABEL": {
          "type": "string"
        },
        "porbps:NAME": {
          "type": "string"
        },
        "porbps:OBJECTID": {
          "type": "string"
        },
        "porbps:SHARED": {
          "type": "string"
        },
        "qs:a0": {
          "type": "string"
        },
        "qs:a0_alt": {
          "type": "string"
        },
        "qs:a0_lc": {
          "type": "string"
        },
        "qs:a1": {
          "type": "string"
        },
        "qs:a1_alt": {
          "type": "string"
        },
        "qs:a1_lc": {
          "type": "string"
        },
        "qs:a1r": {
          "type": "string"
        },
        "qs:a1r_alt": {
          "type": "string"
        },
        "qs:a1r_lc": {
          "type": "string"
        },
        "qs:a2": {
          "type": "string"
        },
        "qs:a2_alt": {
          "type": "string"
        },
        "qs:a2_lc": {
          "type": "string"
        },
        "qs:a2r": {
          "type": "string"
        },
        "qs:a2r_alt": {
          "type": "string"
        },
        "qs:a2r_lc": {
          "type": "string"
        },
        "qs:adm0": {
          "type": "string"
        },
        "qs:adm0_a3": {
          "type": "string"
        },
        "qs:gn_country": {
          "type": [
            "string",
            "null"
          ]
        },
        "qs:gn_fcode": {
          "type": [
            "string",
            "null"
          ]
        },
        "qs:gn_id": {
          "type": "integer"
        },
        "qs:gn_pop": {
          "type": "integer"
        },
        "qs:id": {
          "type": "integer"
        },
        "qs:la": {
          "type": "string"
        },
        "qs:la_alt": {
          "type": "string"
        },
        "qs:la_lc": {
          "type": "string"
        },
        "qs:level": {
          "type": "string"
        },
        "qs:loc": {
          "type": "string"
        },
        "qs:loc_alt": {
          "type": "string"
        },
        "qs:loc_lc": {
          "type": "string"
        },
        "qs:name": {
          "type": "string"
        },
        "qs:name_adm0": {
          "type": "string"
        },
        "qs:name_adm1": {
          "type": [
            "string",
            "null"
          ]
        },
        "qs:photos": {
          "type": "integer"
        },
        "qs:photos_1k": {
          "type": "integer"
        },
        "qs:photos_9k": {
          "type": "integer"
        },
        "qs:photos_9r": {
          "type": "integer"
        },
        "qs:photos_all": {
          "type": "integer"
        },
        "qs:photos_sr": {
          "type": "integer"
        },
        "qs:placetype": {
          "type": "string"
        },
        "qs:pop": {
          "type": "integer"
        },
        "qs:pop_sr": {
          "type": "string"
        },
        "qs:qs_id": {
          "type": "integer"
        },
        "qs:scale": {
          "type": "integer"
        },
        "qs:source": {
          "type": "string"
        },
        "qs:type": {
          "type": "string"
        },
        "qs:woe_adm0": {
          "type": "integer"
        },
        "qs:woe_id": {
          "type": "integer"
        },
        "qs_pg:aaroncc": {
          "type": "string"
        },
        "qs_pg:gn_country": {
          "type": [
            "string",
            "null"
          ]
        },
        "qs_pg:gn_fcode": {
          "type": "string"
        },
        "qs_pg:gn_id": {
          "type": "integer"
        },
        "qs_pg:gn_pop": {
          "type": "integer"
        },
        "qs_pg:name": {
          "type": "string"
        },
        "qs_pg:name_adm0": {
          "type": "string"
        },
        "qs_pg:name_adm1": {
          "type": "string"
        },
        "qs_pg:photos": {
          "type": "integer"
        },
        "qs_pg:photos_1k": {
          "type": "integer"
        },
        "qs_pg:photos_9k": {
          "type": "integer"
        },
        "qs_pg:photos_9r": {
          "type": "integer"
        },
        "qs_pg:photos_all": {
          "type": "integer"
        },
        "qs_pg:photos_sr": {
          "type": "integer"
        },
        "qs_pg:pop_sr": {
          "type": "integer"
        },
        "qs_pg:qs_id": {
          "type": "integer"
        },
        "qs_pg:qs_pg_placetype": {
          "type": "string"
        },
        "qs_pg:qs_pg_placetype_gp": {
          "type": "string"
        },
        "qs_pg:woe_adm0": {
          "type": "integer"
        },
        "qs_pg:woe_id": {
          "type": "integer"
        },
        "resto:alcohol": {
          "type": "string"
        },
        "resto:credit_cards": {
          "type": "string"
        },
        "resto:vegan_friendly": {
          "type": "string"
        },
        "resto:vegetarian_friendly": {
          "type": "string"
        },
        "reversegeo:geometry": {
          "type": "string"
        },
        "reversegeo:latitude": {
          "type": "number"
        },
        "reversegeo:longitude": {
          "type": "number"
        },
        "reversegeo:polygon": {
          "type": "string"
        },
        "santabar:NEIGHBORHOOD": {
          "type": "string"
        },
        "santabar:NHOOD": {
          "type": "string"
        },
        "sdgis:NAME": {
          "type": "string"
        },
        "seagv:hoods_": {
          "type": "string"
        },
        "seagv:hoods_id": {
          "type": "string"
        },
        "seagv:l_hood": {
          "type": "string"
        },
        "seagv:l_hoodid": {
          "type": "string"
        },
        "seagv:objectid": {
          "type": "string"
        },
        "seagv:s_hood": {
          "type": "string"
        },
        "seagv:symbol": {
          "type": "string"
        },
        "seagv:symbol2": {
          "type": "string"
        },
        "sfac:accession_id": {
          "type": "string"
        },
        "sfac:artist": {
          "type": "string"
        },
        "sfac:created_at": {
          "type": "string"
        },
        "sfac:credit_line": {
          "type": "string"
        },
        "sfac:display_dimensions": {
          "type": "string"
        },
        "sfac:id": {
          "type": "string"
        },
        "sfac:location": {
          "type": "string"
        },
        "sfac:location_description": {
          "type": "string"
        },
        "sfac:medium": {
          "type": "string"
        },
        "sfac:revision": {
          "type": "string"
        },
        "sfac:source": {
          "type": "string"
        },
        "sfac:title": {
          "type": "string"
        },
        "sfgov:dept": {
          "type": "string"
        },
        "sfgov:deptname": {
          "type": "string"
        },
        "sfgov:facility_i": {
          "type": "string"
        },
        "sfgov:link": {
          "type": "string"
        },
        "sfgov:school_typ": {
          "type": "string"
        },
        "sg:address": {
          "type": "string"
        },
        "sg:categories": {
          "type": "string"
        },
        "sg:city": {
          "type": "string"
        },
        "sg:classifiers": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "category": {
                "type": "string"
              },
              "type": {
                "type": "string"
              },
              "subcategory": {
                "type": "string"
              }
            }
          }
        },
        "sg:href": {
          "type": "string"
        },
        "sg:menulink": {
          "type": "string"
        },
        "sg:orig_category": {
          "type": "string"
        },
        "sg:owner": {
          "type": "string"
        },
        "sg:phone": {
          "type": "string"
        },
        "sg:postcode": {
          "type": "string"
        },
        "sg:province": {
          "type": "string"
        },
        "sg:tags": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "sg:website": {
          "type": "string"
        },
        "SIJ:admin_1r": {
          "type": "string"
        },
        "SIJ:hasc_id": {
          "type": "string"
        },
        "src:centroid_lbl": {
          "type": "string"
        },
        "src:geom": {
          "type": "string"
        },
        "src:geom_alt": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "src:lbl": {
          "type": "string"
        },
        "src:lbl_centroid": {
          "type": "string"
        },
        "src:population": {
          "type": "string"
        },
        "ssuberlin:gml_id": {
          "type": "string"
        },
        "ssuberlin:spatial_alias": {
          "type": "string"
        },
        "ssuberlin:spatial_name": {
          "type": "string"
        },
        "statoids:area_km": {
          "type": "string"
        },
        "statoids:area_mi": {
          "type": "string"
        },
        "statoids:areakm": {
          "type": "string"
        },
        "statoids:as-of-date": {
          "type": "string"
        },
        "statoids:capital": {
          "type": "string"
        },
        "statoids:country": {
          "type": "string"
        },
        "statoids:date": {
          "type": "string"
        },
        "statoids:dial": {
          "type": "string"
        },
        "statoids:ds": {
          "type": "string"
        },
        "statoids:fifa": {
          "type": "string"
        },
        "statoids:fips": {
          "type": "string"
        },
        "statoids:gaul": {
          "type": "string"
        },
        "statoids:gec": {
          "type": "string"
        },
        "statoids:hasc": {
          "type": "string"
        },
        "statoids:independent": {
          "type": "string"
        },
        "statoids:inegi": {
          "type": "string"
        },
        "statoids:ioc": {
          "type": "string"
        },
        "statoids:iso": {
          "type": "string"
        },
        "statoids:iso_a2": {
          "type": "string"
        },
        "statoids:iso_a3": {
          "type": "string"
        },
        "statoids:iso_num": {
          "type": "string"
        },
        "statoids:itu": {
          "type": "string"
        },
        "statoids:marc": {
          "type": "string"
        },
        "statoids:name": {
          "type": "string"
        },
        "statoids:population": {
          "type": "string"
        },
        "statoids:statoid": {
          "type": "string"
        },
        "statoids:timezone": {
          "type": "string"
        },
        "statoids:type": {
          "type": "string"
        },
        "statoids:tz": {
          "type": "string"
        },
        "statoids:wmo": {
          "type": "string"
        },
        "stpaulgov:district": {
          "type": "string"
        },
        "stpaulgov:name1": {
          "type": "string"
        },
        "stpaulgov:name2": {
          "type": "string"
        },
        "tkugov:ID": {
          "type": "string"
        },
        "tkugov:NIMI": {
          "type": "string"
        },
        "tkugov:PIENALUE": {
          "type": "string"
        },
        "tmpgov:id": {
          "type": "string"
        },
        "tmpgov:MI_PRINX": {
          "type": "string"
        },
        "tmpgov:name": {
          "type": "string"
        },
        "torsdfa:AREA_NAME": {
          "type": "string"
        },
        "torsdfa:AREA_S_CD": {
          "type": "string"
        },
        "tourist:latitude": {
          "type": "number"
        },
        "tourist:longitude": {
          "type": "number"
        },
        "unlc:subdivision": {
          "type": "string"
        },
        "uscensus:aland": {
          "type": "string"
        },
        "uscensus:awater": {
          "type": "string"
        },
        "uscensus:cd115fp": {
          "type": "string"
        },
        "uscensus:cdsessn": {
          "type": "string"
        },
        "uscensus:funcstat": {
          "type": "string"
        },
        "uscensus:geoid": {
          "type": "string"
        },
        "uscensus:intptlat": {
          "type": "string"
        },
        "uscensus:intptlon": {
          "type": "string"
        },
        "uscensus:lsad": {
          "type": "string"
        },
        "uscensus:lsy": {
          "type": "string"
        },
        "uscensus:mtfcc": {
          "type": "string"
        },
        "uscensus:namelsad": {
          "type": "string"
        },
        "uscensus:sldlst": {
          "type": "string"
        },
        "uscensus:sldust": {
          "type": "string"
        },
        "uscensus:statefp": {
          "type": "string"
        },
        "vanpds:MAPID": {
          "type": "string"
        },
        "vanpds:NAME": {
          "type": "string"
        },
        "wapo:quadrant": {
          "type": "string"
        },
        "wapo:subhood": {
          "type": "string"
        },
        "wk:area": {
          "type": "number"
        },
        "wk:elevation": {
          "type": "number"
        },
        "wk:lat": {
          "type": "string"
        },
        "wk:latitude": {
          "type": "number"
        },
        "wk:long": {
          "type": "string"
        },
        "wk:longitude": {
          "type": "number"
        },
        "wk:population": {
          "type": "number"
        },
        "wk:wordcount": {
          "type": "integer"
        },
        "woe:adm0_id": {
          "type": "integer"
        },
        "woe:hierarchy": {
          "type": "object",
          "properties": {
            "state_id": {
              "type": "integer"
            },
            "planet_id": {
              "type": "integer"
            },
            "country_id": {
              "type": "integer"
            },
            "county_id": {
              "type": "integer"
            }
          }
        },
        "woe:name_adm0": {
          "type": "string"
        },
        "woe:name_adm1": {
          "type": "string"
        },
        "woe:placetype": {
          "type": "string"
        },
        "zolk:description": {
          "type": "string"
        },
        "zolk:Name": {
          "type": "string"
        },
        "zs:area_m": {
          "type": "number"
        },
        "zs:blockids": {
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "zs:counties": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "zs:housing10": {
          "type": "integer"
        },
        "zs:id": {
          "type": "integer"
        },
        "zs:label": {
          "type": "string"
        },
        "zs:name": {
          "type": "string"
        },
        "zs:placetype": {
          "type": "string"
        },
        "zs:pop10": {
          "type": "integer"
        }
      },
      "patternProperties": {
        "^abrv:*_x_colloquial$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^abrv:*_x_historical$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^abrv:*_x_preferred$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^abrv:*_x_unknown$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^abrv:*_x_variant$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^name:*_x_colloquial$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^name:*_x_historical$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^name:*_x_preferred$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^name:*_x_unknown$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "^name:*_x_variant$": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    }
  }
}`
