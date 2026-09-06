#include "texturesys_wrapper.h"
#include <OpenImageIO/texture.h>
#include <cstdlib>
#include <cstring>

using namespace OIIO;

struct TextureSystemHandle {
  std::shared_ptr<TextureSystem> ts;
};

namespace {

Tex::Wrap to_oiio_wrap(TextureWrapMode w) {
  switch (w) {
  case WRAP_BLACK:
    return Tex::Wrap::Black;
  case WRAP_CLAMP:
    return Tex::Wrap::Clamp;
  case WRAP_PERIODIC:
    return Tex::Wrap::Periodic;
  case WRAP_MIRROR:
    return Tex::Wrap::Mirror;
  case WRAP_DEFAULT:
  default:
    return Tex::Wrap::Default;
  }
}

TextureOpt to_oiio_opt(const TextureLookupOptions *opts) {
  TextureOpt opt;
  opt.firstchannel = opts->firstchannel;
  opt.subimage = opts->subimage;
  opt.swrap = to_oiio_wrap(opts->swrap);
  opt.twrap = to_oiio_wrap(opts->twrap);
  opt.sblur = opts->sblur;
  opt.tblur = opts->tblur;
  opt.swidth = opts->swidth;
  opt.twidth = opts->twidth;
  opt.fill = opts->fill;
  return opt;
}

} // namespace

TextureSystemHandle *texturesystem_create(char **error_msg) {
  auto ts = TextureSystem::create(/*shared=*/true);
  if (!ts) {
    if (error_msg)
      *error_msg = strdup("Could not create OIIO TextureSystem");
    return nullptr;
  }
  auto *handle = new TextureSystemHandle();
  handle->ts = ts;
  return handle;
}

void texturesystem_destroy(TextureSystemHandle *handle) { delete handle; }

int texturesystem_texture(TextureSystemHandle *handle, const char *filename,
                          const TextureLookupOptions *opts, float s, float t,
                          float dsdx, float dtdx, float dsdy, float dtdy,
                          int nchannels, float *result, char **error_msg) {
  if (!handle || !handle->ts) {
    if (error_msg)
      *error_msg = strdup("Invalid TextureSystem handle");
    return 1;
  }
  TextureOpt opt = to_oiio_opt(opts);
  bool ok = handle->ts->texture(ustring(filename), opt, s, t, dsdx, dtdx, dsdy,
                                dtdy, nchannels, result);
  if (!ok) {
    if (error_msg)
      *error_msg = strdup(handle->ts->geterror().c_str());
    return 1;
  }
  return 0;
}

int texturesystem_get_texture_info_int(TextureSystemHandle *handle,
                                       const char *filename,
                                       const char *dataname, int *out,
                                       char **error_msg) {
  if (!handle || !handle->ts) {
    if (error_msg)
      *error_msg = strdup("Invalid TextureSystem handle");
    return 1;
  }
  bool ok = handle->ts->get_texture_info(ustring(filename), 0,
                                         ustring(dataname), TypeDesc::INT, out);
  if (!ok) {
    if (error_msg)
      *error_msg = strdup(handle->ts->geterror().c_str());
    return 1;
  }
  return 0;
}

int texturesystem_get_texture_info_float(TextureSystemHandle *handle,
                                         const char *filename,
                                         const char *dataname, float *out,
                                         char **error_msg) {
  if (!handle || !handle->ts) {
    if (error_msg)
      *error_msg = strdup("Invalid TextureSystem handle");
    return 1;
  }
  bool ok = handle->ts->get_texture_info(
      ustring(filename), 0, ustring(dataname), TypeDesc::FLOAT, out);
  if (!ok) {
    if (error_msg)
      *error_msg = strdup(handle->ts->geterror().c_str());
    return 1;
  }
  return 0;
}
