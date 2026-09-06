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
  float sblur, tblur;    // Blur amount
  float swidth, twidth;  // Multiplier for derivatives (OIIO default: 1)
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

#ifdef __cplusplus
}
#endif
