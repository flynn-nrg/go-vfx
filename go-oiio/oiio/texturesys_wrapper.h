#pragma once

#ifdef __cplusplus
extern "C" {
#endif

typedef struct TextureSystemHandle TextureSystemHandle;

typedef enum {
  WRAP_DEFAULT = 0,
  WRAP_BLACK,
  WRAP_CLAMP,
  WRAP_PERIODIC,
  WRAP_MIRROR,
} TextureWrapMode;

typedef struct TextureLookupOptions {
  int firstchannel;      // First channel of the lookup
  int subimage;          // Subimage/face index
  TextureWrapMode swrap;
  TextureWrapMode twrap;
  TextureWrapMode rwrap;         // Wrap mode in the r (volumetric) direction
  float sblur, tblur, rblur;     // Blur amount
  float swidth, twidth, rwidth;  // Multiplier for derivatives (OIIO default: 1)
  float fill;            // Fill value for channels beyond what the file has
} TextureLookupOptions;

// Creates a shared OIIO::TextureSystem (owns its own tile cache). The handle
// must be released with texturesystem_destroy. Returns NULL on failure and
// sets *error_msg (caller must free() it).
TextureSystemHandle *texturesystem_create(char **error_msg);
void texturesystem_destroy(TextureSystemHandle *ts);

// Filtered texture lookup at (s,t) using screen-space derivatives dsdx/dtdx
// (d/dx) and dsdy/dtdy (d/dy) to select mip level and filter width, mirroring
// OSL's texture() shadeop. result must have room for nchannels floats.
// Returns 0 on success, non-zero on failure (error message via error_msg,
// caller must free() it).
int texturesystem_texture(TextureSystemHandle *ts, const char *filename,
                          const TextureLookupOptions *opts, float s, float t,
                          float dsdx, float dtdx, float dsdy, float dtdy,
                          int nchannels, float *result, char **error_msg);

// Filtered 3D (volumetric) texture lookup at point p, using the screen-space
// derivatives dpdx/dpdy/dpdz (each a float[3]) to select filter width,
// mirroring OSL's texture3d() shadeop. p/dpdx/dpdy/dpdz are each float[3]
// (x,y,z). result must have room for nchannels floats. Returns 0 on success,
// non-zero on failure (error message via error_msg, caller must free() it).
int texturesystem_texture3d(TextureSystemHandle *ts, const char *filename,
                            const TextureLookupOptions *opts, const float *p,
                            const float *dpdx, const float *dpdy,
                            const float *dpdz, int nchannels, float *result,
                            char **error_msg);

// Filtered environment lookup along direction r, using the screen-space
// derivatives drdx/drdy (each a float[3]) to select filter width, mirroring
// OSL's environment() shadeop. r/drdx/drdy are each float[3] (x,y,z). result
// must have room for nchannels floats. Returns 0 on success, non-zero on
// failure (error message via error_msg, caller must free() it).
int texturesystem_environment(TextureSystemHandle *ts, const char *filename,
                              const TextureLookupOptions *opts,
                              const float *r, const float *drdx,
                              const float *drdy, int nchannels, float *result,
                              char **error_msg);

// Metadata queries mirroring OSL's gettextureinfo() shadeop, e.g. dataname
// "exists", "channels", "resolution" (int[2]), "worldtocamera" (float[16]).
int texturesystem_get_texture_info_int(TextureSystemHandle *ts,
                                       const char *filename,
                                       const char *dataname, int *out,
                                       char **error_msg);
int texturesystem_get_texture_info_float(TextureSystemHandle *ts,
                                         const char *filename,
                                         const char *dataname, float *out,
                                         char **error_msg);

// TextureSystem-wide configuration, e.g. name "searchpath" (string),
// "max_open_files"/"autotile"/"automip" (int), "max_memory_MB" (float).
// Returns 0 on success, non-zero if the attribute name/type is not
// recognized (error message via error_msg, caller must free() it).
int texturesystem_attribute_int(TextureSystemHandle *ts, const char *name,
                                int value, char **error_msg);
int texturesystem_attribute_float(TextureSystemHandle *ts, const char *name,
                                  float value, char **error_msg);
int texturesystem_attribute_string(TextureSystemHandle *ts, const char *name,
                                   const char *value, char **error_msg);

// Reads back TextureSystem-wide configuration set via texturesystem_attribute_*.
// texturesystem_getattribute_string writes a newly allocated string to *out;
// the caller must free() it. Returns 0 on success, non-zero if the attribute
// name/type is not recognized (error message via error_msg, caller must
// free() it).
int texturesystem_getattribute_int(TextureSystemHandle *ts, const char *name,
                                   int *out, char **error_msg);
int texturesystem_getattribute_float(TextureSystemHandle *ts, const char *name,
                                     float *out, char **error_msg);
int texturesystem_getattribute_string(TextureSystemHandle *ts,
                                      const char *name, char **out,
                                      char **error_msg);

#ifdef __cplusplus
}
#endif
